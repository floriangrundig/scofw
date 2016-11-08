module Update exposing (..)

import Messages exposing (Msg(..))
import Models exposing (Model)
import LiveView.Messages as LiveViewMessage exposing (Msg(..))
import LiveView.Update as LiveView exposing (update)
import ServerMsgDecoder


update : Messages.Msg -> Model -> ( Model, Cmd Messages.Msg )
update msg model =
    case msg of
        LiveViewMsg subMsg ->
            let
                ( updatedLiveViewModel, cmd ) =
                    LiveView.update subMsg model.liveViewModel
            in
                ( { model | liveViewModel = updatedLiveViewModel }, Cmd.map LiveViewMsg cmd )

        WsUpdate str ->
            let
                serverMsg =
                    str
                        |> ServerMsgDecoder.decodeServerMsg
                        |> Debug.log "Incoming Websocket Msg:"

                ( updatedLiveViewModel, cmd ) =
                    case serverMsg of
                        Err m ->
                            let
                                _ =
                                    Debug.log "Error parsing server msg" m
                            in
                                ( model.liveViewModel, Cmd.none )

                        Ok m ->
                            LiveView.update (LiveViewMessage.FileChanged m) model.liveViewModel
            in
                ( { model | liveViewModel = updatedLiveViewModel }, Cmd.map LiveViewMsg cmd )


