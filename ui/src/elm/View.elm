module View exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html
import Messages exposing (Msg(..))
import Models exposing (Model)
import LiveView.View as LiveView exposing (view)
import Routing exposing (Route(..))


view : Model -> Html Msg
view model =
    div []
        [ page model ]


page : Model -> Html Msg
page model =
    case model.route of
        RootRoute ->
            rootView

        LiveViewRoute ->
            Html.map LiveViewMsg (LiveView.view model.liveViewModel)

        NotFoundRoute ->
            notFoundView


rootView : Html msg
rootView =
    div []
        [ text "Overview"
        , div []
            [ a [ href "#live" ] [ text "Live View" ] ]
        ]


notFoundView : Html msg
notFoundView =
    div []
        [ text "Not found"
        ]
