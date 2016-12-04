module CommonModels exposing (..)


type alias File =
    { name : String
    }


type alias Diff =
    { file : File
    , diff : String
    }


type alias Session =
    { name : String
    , diffs : List Diff
    }


type alias FileChangeEvent =
    { filename : String
    , op : Op
    }


type alias ServerMsg =
    { fileChanges : List (Result String FileChange)
    , currentSession : String
    , projectName : String
    }


type alias FileChange =
    { oldName : Result String String
    , newName : Result String String
    , op : Op
    , hunks : List (Result String Hunk)
    , time : Result String String
    }


type alias Hunk =
    { context : Maybe String
    , ranges : Result String HunkRanges
    , lines : List HunkLine
    , additions : Int
    , deletions : Int
    }


type alias HunkRanges =
    { fromFileLineNumberStart : Int
    , fromFileLineNumberEnd : Int
    , toFileLineNumberStart : Int
    , toFileLineNumberEnd : Int
    }


type HunkLine
    = Context String
    | Addition String
    | Deletion String


type Op
    = Modified
    | Added
    | Removed
    | Renamed
