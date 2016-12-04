module Main exposing (..)

import Navigation
import Models exposing (Model, initialModel)
import View exposing (view)
import Update exposing (update)
import Messages exposing (Msg(..))
import Routing exposing (Route)
import WebSocket


-- INITIALIZATION


init : Navigation.Location -> ( Model, Cmd Msg )
init location =
    let
        currentRoute =
            Routing.parseHash location
    in
        ( initialModel currentRoute, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    WebSocket.listen "ws://localhost:5000/ws" WsUpdate



-- MAIN


main : Program Never Model Msg
main =
    Navigation.program UrlChange
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
