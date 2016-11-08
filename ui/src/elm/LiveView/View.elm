module LiveView.View exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import LiveView.Messages exposing (Msg(..))
import LiveView.Models exposing (..)
import CommonModels as Common exposing (..)
import String

view : Model -> Html Msg
view model =
    div []
        [ pageDom model ]


pageDom : Model -> Html Msg
pageDom model =
    div [ class "liveview__page" ]
        [ case model.state of
            Initial ->
                div [] [ text "LiveView waiting for events" ]

            ReceivingEvents s ->
                (eventsDom s.serverMessages)
        ]


eventsDom : List Common.ServerMsg -> Html Msg
eventsDom serverMsgs =
    div [ class "events__container" ]
        [ div [ class "events__summary" ] (eventsHeaderDom serverMsgs)
        , div [ class "events__listing" ] [ eventsListingDom serverMsgs ]
        ]


eventsListingDom : List Common.ServerMsg -> Html Msg
eventsListingDom serverMsgs =
    let
        numberOfEvents =
            List.length serverMsgs
    in
        div [ class "" ]
            [ ul []
                (List.map eventDom serverMsgs)
            ]


eventDom : Common.ServerMsg -> Html Msg
eventDom msg =
    let
        patch =
            msg.patch |> Debug.log "events"

        patchItems = String.split "\n" patch
                        |> Debug.log "formatted"
    in
        li [ class "events__item code" ]
            [pre [] [text patch]]
            


eventsHeaderDom : List Common.ServerMsg -> List (Html Msg)
eventsHeaderDom serverMsgs =
    let
        numberOfEvents =
            List.length serverMsgs
    in
        [ div [ class "table-row table-header" ]
            [ div [] [ text "#Modifications" ]
            ]
        , div [ class "table-row" ]
            [ text <| toString numberOfEvents ]
        ]
