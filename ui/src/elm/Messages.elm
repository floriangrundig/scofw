module Messages exposing (..)

import LiveView.Messages


type Msg
    = LiveViewMsg LiveView.Messages.Msg
    | WsUpdate String
