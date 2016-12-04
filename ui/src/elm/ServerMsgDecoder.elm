module ServerMsgDecoder exposing (..)

import Json.Decode as Json exposing (..)
import CommonModels exposing (..)
import Parser.UnifiedDiffParser exposing (parse)


decodeServerMsg : String -> Result String ServerMsg
decodeServerMsg =
    Json.decodeString
        (Json.map3 ServerMsg
            (field "FileChanges" fileChangesDecoder)
            (field "CurrentScoSession" Json.string)
            (field "ProjectName" Json.string)
        )


fileChangesDecoder : Json.Decoder (List (Result String FileChange))
fileChangesDecoder =
    Json.andThen
        (\logEntry -> Json.succeed (parse logEntry))
        Json.string
