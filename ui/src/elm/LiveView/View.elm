module LiveView.View exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import LiveView.Messages exposing (Msg(..))
import LiveView.Models exposing (..)
import CommonModels as Common exposing (..)


view : Model -> Html Msg
view model =
    div []
        [ pageDom model ]


pageDom : Model -> Html Msg
pageDom model =
    div [ class "liveview__page" ]
        [ case model.state of
            Initial ->
                div [] [ text "LiveView waitig for events" ]

            ReceivingEvents s ->
                (liveViewDom s.serverMessages)
        ]


liveViewDom : List Common.ServerMsg -> Html Msg
liveViewDom serverMsgs =
    div [ class "events__container" ]
        [ div [ class "events__summary" ] (eventsHeaderDom serverMsgs)
        , div [ class "events__listing" ]
            [ div [ class "events__modification-graph" ]
                [ eventsModificationGraph serverMsgs ]
            , div [ class "events__details" ]
                (eventsListingDom serverMsgs)
            ]
        ]


eventsListingDom : List Common.ServerMsg -> List (Html Msg)
eventsListingDom serverMsgs =
    let
        numberOfEvents =
            List.length serverMsgs

        fileChanges =
            List.concat (List.map (\m -> List.map eventResultDom m.fileChanges) serverMsgs)
    in
        fileChanges


eventResultDom : Result String FileChange -> Html Msg
eventResultDom result =
    let
        _ =
            Debug.log "eventDom" result
    in
        case result of
            Ok r ->
                eventResultOkDom r

            Err m ->
                eventResultErrDom m


eventResultOkDom : FileChange -> Html Msg
eventResultOkDom fileChange =
    div [ class "event" ]
        [ eventSummaryDom fileChange
        , eventDiffDom fileChange
        ]


eventSummaryDom : FileChange -> Html Msg
eventSummaryDom fileChange =
    let
        fileName =
            case fileChange.newName of
                Ok f ->
                    f

                Err m ->
                    m

        time =
            case fileChange.time of
                Ok t ->
                    t

                Err m ->
                    m
    in
        div [ class "event__summary" ]
            [ div [] [ text fileName ]
            , div [] [ text time ]
            ]


eventDiffDom : FileChange -> Html Msg
eventDiffDom fileChange =
    div [ class "event__diff" ]
        (List.map hunkResultDom fileChange.hunks)


hunkResultDom : Result String Hunk -> Html Msg
hunkResultDom hunk =
    case hunk of
        Ok h ->
            hunkResultOkDom h

        Err m ->
            hunkResultErrDom m


hunkResultOkDom : Hunk -> Html Msg
hunkResultOkDom hunk =
    div [ class "hunk" ]
        [ hunkHeaderDom hunk
        , hunkDetailsDom hunk
        ]


hunkHeaderDom : Hunk -> Html Msg
hunkHeaderDom hunk =
    div [ class "hunk__summary" ]
        [ text ("Context: " ++ (Maybe.withDefault "-" hunk.context))
        ]


hunkDetailsDom : Hunk -> Html Msg
hunkDetailsDom hunk =
    div [ class "hunk__details" ]
        [ hunkDetailsAfterChangeDom hunk
        , hunkDetailsBeforChangeDom hunk
        ]


hunkDetailCodeLineDom : Bool -> HunkLine -> Html Msg
hunkDetailCodeLineDom isForAfterChange hl =
    case ( isForAfterChange, hl ) of
        ( _, Context l ) ->
            div [ class "code code__context" ] [ text l ]

        ( True, Addition l ) ->
            div [ class "code code__addition" ] [ text l ]

        ( False, Addition l ) ->
            div [ class "code" ] [ text "" ]

        ( True, Deletion l ) ->
            div [ class "code" ] [ text "" ]

        ( False, Deletion l ) ->
            div [ class "code code__deletion" ] [ text l ]


hunkDetailsAfterChangeDom : Hunk -> Html Msg
hunkDetailsAfterChangeDom hunk =
    div [ class "hunk__details-content hunk__details-after" ]
        [ div [] []
        , (div [ class "hunk__details-code" ]
            (List.map (hunkDetailCodeLineDom True) hunk.lines)
          )
        ]


hunkDetailsBeforChangeDom : Hunk -> Html Msg
hunkDetailsBeforChangeDom hunk =
    div [ class "hunk__details-content hunk__details-before" ]
        [ div [] []
        , (div [ class "hunk__details-code" ]
            (List.map (hunkDetailCodeLineDom False) hunk.lines)
          )
        ]


hunkResultErrDom : String -> Html Msg
hunkResultErrDom msg =
    div [ classList [ ( "hunk", True ), ( "parse-error", True ) ] ]
        [ text msg ]


eventResultErrDom : String -> Html Msg
eventResultErrDom msg =
    div [ classList [ ( "event", True ), ( "parse-error", True ) ] ]
        [ text msg ]


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


eventsModificationGraph : List Common.ServerMsg -> Html Msg
eventsModificationGraph msgs =
    text ""
