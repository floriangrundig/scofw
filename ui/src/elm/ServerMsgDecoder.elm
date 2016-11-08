module ServerMsgDecoder exposing (..)

import Json.Decode as Json exposing (..)
import CommonModels exposing (..)
import Parser.UnifiedDiffParser exposing (parse)


decodeServerMsg : String -> Result String ServerMsg
decodeServerMsg =

    Json.decodeString
        (Json.object3 ServerMsg
            ("FileChanges" := fileChangesDecoder)
            ("CurrentScoSession" := Json.string)
            ("ProjectName" := Json.string)
        )


fileChangesDecoder : Json.Decoder (List (Result String FileChange))
fileChangesDecoder =
    Json.customDecoder Json.string (\logEntry -> Ok (parse logEntry))
