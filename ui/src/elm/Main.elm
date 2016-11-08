module Main exposing (..)

import Navigation
import Models exposing (Model, initialModel)
import View exposing (view)
import Update exposing (update)
import Messages exposing (Msg(..))
import Routing exposing (Route)
import WebSocket


init : Result String Route -> ( Model, Cmd Msg )
init result =
    let
        currentRoute =
            Routing.routeFromResult result
    in
        ( initialModel currentRoute, Cmd.none )


urlUpdate : Result String Route -> Model -> ( Model, Cmd Msg )
urlUpdate result model =
    let
        currentRoute =
            Routing.routeFromResult result
    in
        ( { model | route = currentRoute }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    WebSocket.listen "ws://localhost:5000/ws" WsUpdate 



-- MAIN


main : Program Never
main =
    Navigation.program Routing.parser
        { init = init
        , view = view
        , urlUpdate = urlUpdate
        , update = update
        , subscriptions = subscriptions
        }
