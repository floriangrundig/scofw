module Parser.UnifiedDiffParser exposing (parse)

import String
import Result
import Regex exposing (..)
import CommonModels exposing (..)


parse : String -> List (Result String FileChange)
parse str =
    let
        _ =
            Debug.log "parse  entry" str

        logEntries =
            splitLogEntries str
    in
        List.map parseLogEntry logEntries


{-| This is a regex which is able to split up a whole string into all log events
-}
entryStartRegEx : Regex
entryStartRegEx =
    Regex.regex "^\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} diff --git a/.+ b/.+$"


splitLogEntries : String -> List String
splitLogEntries str =
    String.split "\n" str
        |> List.foldr
            (\line acc ->
                if Regex.contains entryStartRegEx line then
                    case acc of
                        [] ->
                            [ line ]

                        x :: xs ->
                            "" :: (line ++ "\n" ++ x) :: xs
                else
                    case acc of
                        [] ->
                            [ line ]

                        x :: xs ->
                            (line ++ "\n" ++ x) :: xs
            )
            []
        |> List.tail
        |> Maybe.withDefault []


parseLogEntry : String -> Result String FileChange
parseLogEntry entry =
    let
        oldName =
            parseOldFileName entry

        newName =
            parseNewFileName entry

        op =
            case ( oldName, newName ) of
                ( Ok "/dev/null", _ ) ->
                    Added

                ( _, Ok "/dev/null" ) ->
                    Removed

                _ ->
                    Modified

        time =
            parseLogTime entry

        hunks =
            parseHunks entry
    in
        Ok
            { oldName = oldName
            , newName = newName
            , op = op
            , time = time
            , hunks = hunks
            }


parseHunks : String -> List (Result String Hunk)
parseHunks entry =
    let
        entryAsLines =
            String.split "\n" entry

        hunkStartRegex =
            Regex.regex "^@@ -\\d+,?\\d* \\+\\d+,?\\d* @@.*"

        hunkStartIndicies =
            List.foldr
                (\( idx, line ) acc ->
                    if Regex.contains hunkStartRegex line then
                        idx :: acc
                    else
                        acc
                )
                []
                (List.indexedMap (,) entryAsLines)
    in
        List.foldr
            (\idx acc ->
                case List.head acc of
                    Nothing ->
                        ( idx, range idx ((List.length entryAsLines) - 1) entryAsLines ) :: acc

                    Just ( lastIdx, _ ) ->
                        ( idx, range idx (lastIdx - 1) entryAsLines ) :: acc
            )
            []
            hunkStartIndicies
            |> List.map (snd >> parseHunk)


parseHunk : List String -> Result String Hunk
parseHunk lines =
    case lines of
        [] ->
            Err "Empty Hunk"

        hunkDetails :: hunkBody ->
            let
                hunkDetailsRegex =
                    Regex.regex "^@@ -(\\d+),?(\\d*) \\+(\\d+),?(\\d*) @@ *(.*)"

                hunkDetailMatches =
                    Regex.find (AtMost 1) hunkDetailsRegex hunkDetails
            in
                case List.head hunkDetailMatches of
                    Nothing ->
                        Err ("Couldn't parse hunk details: " ++ hunkDetails)

                    Just matches ->
                        case matches.submatches of
                            [ fromFileLineNumberStart, fromFileLineNumberEnd, toFileLineNumberStart, toFileLineNumberEnd, context ] ->
                                let
                                    fromFileLineNumberStart' =
                                        fromFileLineNumberStart |> Maybe.withDefault "failed" |> String.toInt

                                    toFileLineNumberStart' =
                                        toFileLineNumberStart |> Maybe.withDefault "failed" |> String.toInt

                                    fromFileLineNumberEnd' =
                                        if fromFileLineNumberEnd == Just "" then
                                            fromFileLineNumberStart'
                                        else
                                            fromFileLineNumberEnd |> Maybe.withDefault "failed" |> String.toInt

                                    toFileLineNumberEnd' =
                                        if toFileLineNumberEnd == Just "" then
                                            toFileLineNumberStart'
                                        else
                                            toFileLineNumberEnd |> Maybe.withDefault "failed" |> String.toInt

                                    ranges =
                                        Result.map4
                                            (\fromStart fromEnd toStart toEnd ->
                                                { fromFileLineNumberStart = fromStart
                                                , toFileLineNumberStart = toStart
                                                , fromFileLineNumberEnd = fromEnd
                                                , toFileLineNumberEnd = toEnd
                                                }
                                            )
                                            fromFileLineNumberStart'
                                            fromFileLineNumberEnd'
                                            toFileLineNumberStart'
                                            toFileLineNumberEnd'

                                    ( additions, deletions, hunkLines ) =
                                        parseHunkBody hunkBody
                                in
                                    Ok
                                        { context = context
                                        , lines = hunkLines
                                        , ranges = ranges
                                        , additions = additions
                                        , deletions = deletions
                                        }

                            -- Err "bae"
                            _ ->
                                Err ("Couldn't parse hunk details: " ++ hunkDetails)


parseHunkBody : List String -> ( Int, Int, List HunkLine )
parseHunkBody lines =
    List.foldr
        (\line ( additions, deletions, lines ) ->
            let
                lineSplit =
                    if line == "" then
                        Just ( ' ', "" )
                        -- test shows that we're reading whitespace lines as empty lines -> that might lead to a trailing context "" for the last hunk
                    else
                        String.uncons line
            in
                case lineSplit of
                    Nothing ->
                        let
                            _ =
                                Debug.log "Error parsing line" line
                        in
                            ( additions, deletions, lines )

                    Just ( h, t ) ->
                        if h == '+' then
                            ( additions + 1, deletions, Addition t :: lines )
                        else if h == '-' then
                            ( additions, deletions + 1, Deletion t :: lines )
                        else if h == ' ' then
                            ( additions, deletions, Context t :: lines )
                        else
                            let
                                _ =
                                    Debug.log "Error parsing line" line
                            in
                                ( additions, deletions, Context t :: lines )
        )
        ( 0, 0, [] )
        lines


range : Int -> Int -> List String -> List String
range from to items =
    (List.indexedMap (,) items)
        |> List.filter (\( idx, line ) -> idx >= from && idx <= to)
        |> List.map snd


parseOldFileName : String -> Result String String
parseOldFileName entry =
    let
        oldNameRegex : Regex
        oldNameRegex =
            Regex.regex "--- a/(.+)\n"
    in
        parseExactlyOneGroup entry "Couldn't parse old file name" oldNameRegex


parseNewFileName : String -> Result String String
parseNewFileName entry =
    let
        newNameRegex : Regex
        newNameRegex =
            Regex.regex "\\+\\+\\+ b/(.+)\n"
    in
        parseExactlyOneGroup entry "Couldn't parse new file name" newNameRegex


parseLogTime : String -> Result String String
parseLogTime entry =
    let
        timeRegex =
            Regex.regex "^(\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}) diff --git a/.+ b/.+\n.*"
    in
        parseExactlyOneGroup entry "Couldn't parse new file name" timeRegex


parseExactlyOneGroup : String -> String -> Regex -> Result String String
parseExactlyOneGroup entry errMsg regex =
    let
        match =
            Regex.find All regex entry

        parseMatch : Regex.Match -> Result String String
        parseMatch m =
            case m.submatches of
                [] ->
                    Err errMsg

                x :: [] ->
                    case x of
                        Just m' ->
                            Ok m'

                        Nothing ->
                            Err errMsg

                _ ->
                    Err errMsg
    in
        case match of
            [] ->
                Err errMsg

            x :: [] ->
                parseMatch x

            _ ->
                Err errMsg
