module Messages exposing (..)

import LiveView.Messages
import Navigation


type Msg
    = LiveViewMsg LiveView.Messages.Msg
    | WsUpdate String
    | UrlChange Navigation.Location
