module LiveView.Models exposing (..)

import CommonModels as Common


type State
    = Initial
    | ReceivingEvents ModelReceiving


type alias Model =
    { state : State
    }


type alias ModelReceiving =
    { serverMessages : List Common.ServerMsg
    }


initialModel : Model
initialModel =
    { state = Initial
    }
