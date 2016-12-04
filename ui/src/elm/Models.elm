module Models exposing (..)

import CommonModels as Common
import LiveView.Models as LiveView exposing (Model)
import Routing exposing (..)


type alias Model =
    { sessions : List Common.Session
    , route : Route
    , liveViewModel : LiveView.Model
    }


initialModel : Route -> Model
initialModel route =
    { sessions = []
    , route = route
    , liveViewModel = LiveView.initialModel
    }
