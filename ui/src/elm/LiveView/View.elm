module LiveView.View exposing (..)

import CommonModels as Common exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import LiveView.Messages exposing (Msg(..))
import LiveView.Models exposing (..)
import Result
import Tuple exposing (..)
import Regex

view : Model -> Html Msg
view model =
    div []
        [ pageDom model ]


pageDom : Model -> Html Msg
pageDom model =
    div [ class "liveview__page" ]
        [ case model.state of
            Initial ->
                div [] [ text "LiveView waitig for events" ]

            ReceivingEvents s ->
                (liveViewDom s.serverMessages)
        ]


liveViewDom : List Common.ServerMsg -> Html Msg
liveViewDom serverMsgs =
    div [ class "events__container" ]
        [ div [ class "events__summary" ] (eventsHeaderDom serverMsgs)
        , div [ class "events__listing" ]
            [ div [ class "events__modification-graph" ]
                [ eventsModificationGraph serverMsgs ]
            , div [ class "events__details" ]
                (eventsListingDom serverMsgs)
            ]
        ]


eventsListingDom : List Common.ServerMsg -> List (Html Msg)
eventsListingDom serverMsgs =
    let
        numberOfEvents =
            List.length serverMsgs

        fileChanges =
            List.concat (List.map (\m -> List.map eventResultDom m.fileChanges) serverMsgs)
    in
        fileChanges


eventResultDom : Result String FileChange -> Html Msg
eventResultDom result =
    let
        _ =
            Debug.log "eventDom" result
    in
        case result of
            Ok r ->
                eventResultOkDom r

            Err m ->
                eventResultErrDom m


eventResultOkDom : FileChange -> Html Msg
eventResultOkDom fileChange =
    div [ class "event" ]
        [ eventSummaryDom fileChange
        , eventDiffDom fileChange
        ]


eventSummaryDom : FileChange -> Html Msg
eventSummaryDom fileChange =
    let
        fileName =
            case fileChange.newName of
                Ok f ->
                    f

                Err m ->
                    m

        time =
            case fileChange.time of
                Ok t ->
                    t

                Err m ->
                    m
    in
        div [ class "event__summary" ]
            [ div [] [ text fileName ]
            , div [] [ text time ]
            ]


eventDiffDom : FileChange -> Html Msg
eventDiffDom fileChange =
    div [ class "event__diff" ]
        (List.map hunkResultDom fileChange.hunks)


hunkResultDom : Result String Hunk -> Html Msg
hunkResultDom hunk =
    case hunk of
        Ok h ->
            hunkResultOkDom h

        Err m ->
            hunkResultErrDom m


hunkResultOkDom : Hunk -> Html Msg
hunkResultOkDom hunk =
    div [ class "hunk" ]
        [ hunkHeaderDom hunk
        , hunkDetailsDom hunk
        ]


hunkHeaderDom : Hunk -> Html Msg
hunkHeaderDom hunk =
    div [ class "hunk__summary" ]
        [ text ("Context: " ++ (Maybe.withDefault "-" hunk.context))
        ]


type ViewableHunkLine
    = Context_ ( Int, String )
    | Addition_ ( Int, String )
    | AdditionPadding
    | Deletion_ ( Int, String )
    | DeletionPadding


appendAdditionPadding : Int -> List ViewableHunkLine -> List ViewableHunkLine
appendAdditionPadding times viewableHunkLines =
    if (times > 0) then
        appendAdditionPadding (times - 1) (AdditionPadding :: viewableHunkLines)
    else
        viewableHunkLines


appendDeletionPadding : Int -> List ViewableHunkLine -> List ViewableHunkLine
appendDeletionPadding times viewableHunkLines =
    if (times > 0) then
        appendDeletionPadding (times - 1) (DeletionPadding :: viewableHunkLines)
    else
        viewableHunkLines


asViewableHunkLines : Hunk -> ( List ViewableHunkLine, List ViewableHunkLine )
asViewableHunkLines hunk =
    let
        helperFn : List HunkLine -> ( Int, Int ) -> ( Int, Int ) -> ( List ViewableHunkLine, List ViewableHunkLine ) -> ( List ViewableHunkLine, List ViewableHunkLine )
        helperFn hunkLines ( lnA, lnD ) ( countA, countD ) (( viewableA, viewableD ) as acc) =
            case hunkLines of
                (Addition lA) :: xs ->
                    let
                        adDiff =
                            countA + 1 - countD

                        viewableA_ =
                            Addition_ ( lnA, lA ) :: viewableA
                    in
                        helperFn xs ( lnA + 1, lnD ) ( countA + 1, countD ) ( viewableA_, viewableD )

                (Deletion lD) :: xs ->
                    let
                        diff =
                            countD + 1 - countA

                        viewableD_ =
                            Deletion_ ( lnD, lD ) :: viewableD
                    in
                        helperFn xs ( lnA, lnD + 1 ) ( countA, countD + 1 ) ( viewableA, viewableD_ )

                (Context lC) :: xs ->
                    let
                        adDiff =
                            countA - countD

                                -- 188
                        viewableD_ =
                            if adDiff > 0 then
                                appendDeletionPadding adDiff viewableD
                            else
                                viewableD

                        viewableA_ =
                            if adDiff < 0 then
                                appendAdditionPadding (adDiff * -1) viewableA
                            else
                                viewableA
                    in
                        helperFn xs ( lnA + 1, lnD + 1 ) ( 0, 0 ) ( Context_ (lnA, lC) :: viewableA_, Context_ (lnD, lC) :: viewableD_ )

                [] ->
                    let
                        adDiff =
                            countA - countD

                        viewableD_ =
                            if adDiff > 0 then
                                appendDeletionPadding adDiff viewableD
                            else
                                viewableD

                        viewableA_ =
                            if adDiff < 0 then
                                appendAdditionPadding (adDiff * -1) viewableA
                            else
                                viewableA
                    in
                        ( List.reverse viewableA_, List.reverse viewableD_ )

        lnA =
            Result.withDefault 0 <| Result.map .fromFileLineNumberStart hunk.ranges

        lnB =
            Result.withDefault 0 <| Result.map .toFileLineNumberStart hunk.ranges
    in
        helperFn hunk.lines ( lnA, lnB ) ( 0, 0 ) ( [], [] )


hunkDetailsDom : Hunk -> Html Msg
hunkDetailsDom hunk =
    let
        viewableHunkLines =
            asViewableHunkLines hunk
    in
        div [ class "hunk__details" ]
            [ hunkDetailsAfterChangeDom <| first viewableHunkLines
            , hunkDetailsBeforChangeDom <| second viewableHunkLines
            ]
            --s
toText : String -> Html Msg
toText line =
    if Regex.contains (Regex.regex "^\\s*$") line then
        br [] []
    else
        text line

hunkDetailCodeLineDom : Bool -> ViewableHunkLine -> Html Msg
hunkDetailCodeLineDom isForAfterChange hl =
    case ( isForAfterChange, hl ) of
        ( _, Context_ ( _, l ) ) ->
            div [ class "code code__context" ] [ toText l ]

        ( True, Addition_ ( _, l ) ) ->
            div [ class "code code__addition" ] [ toText l ]

        ( False, Addition_ ( _, l ) ) ->
            div [ class "code" ] [ text "" ]

        ( True, Deletion_ ( _, l ) ) ->
            div [ class "code" ] [ text "" ]

        ( False, Deletion_ ( _, l ) ) ->
            div [ class "code code__deletion" ] [ toText l ]

        ( False, AdditionPadding ) ->
            text ""

        ( True, AdditionPadding ) ->
            div [ class "code code__padding" ] [ br [] [] ]

        ( True, DeletionPadding ) ->
            text ""

        ( False, DeletionPadding ) ->
            div [ class "code code__padding" ] [ br [] [] ]


hunkDetailsAfterChangeDom : List ViewableHunkLine -> Html Msg
hunkDetailsAfterChangeDom hunkLines =
    div [ class "hunk__details-content hunk__details-after" ]
        [ div [class "hunk__details-lineNumbers"] (List.map hunkDetailLineNumber hunkLines)
        , (div [ class "hunk__details-code" ]
            (List.map (hunkDetailCodeLineDom True) hunkLines)
          )
        ]


hunkDetailsBeforChangeDom : List ViewableHunkLine -> Html Msg
hunkDetailsBeforChangeDom hunkLines =
    div [ class "hunk__details-content hunk__details-before" ]
        [ div [class "hunk__details-lineNumbers"] (List.map hunkDetailLineNumber hunkLines)
        , (div [ class "hunk__details-code" ]
            (List.map (hunkDetailCodeLineDom False) hunkLines)
          )
        ]

hunkDetailLineNumber: ViewableHunkLine -> Html Msg
hunkDetailLineNumber hunkLine =
    case hunkLine of
        Context_ (ln, _) ->
            div [] [text <| toString ln]
        Addition_ (ln, _) ->
                div [] [text <| toString ln]
        Deletion_ (ln, _) ->
                div [] [text <| toString ln]
        AdditionPadding ->
                br [] []
        DeletionPadding ->
            br [] []



hunkResultErrDom : String -> Html Msg
hunkResultErrDom msg =
    div [ classList [ ( "hunk", True ), ( "parse-error", True ) ] ]
        [ text msg ]


eventResultErrDom : String -> Html Msg
eventResultErrDom msg =
    div [ classList [ ( "event", True ), ( "parse-error", True ) ] ]
        [ text msg ]


eventsHeaderDom : List Common.ServerMsg -> List (Html Msg)
eventsHeaderDom serverMsgs =
    let
        numberOfEvents =
            List.length serverMsgs
    in
        [ div [ class "table-row table-header" ]
            [ div [] [ text "#Modifications" ]
            ]
        , div [ class "table-row" ]
            [ text <| toString numberOfEvents ]
        ]


eventsModificationGraph : List Common.ServerMsg -> Html Msg
eventsModificationGraph msgs =
    text ""
