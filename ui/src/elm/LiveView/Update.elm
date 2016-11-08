module LiveView.Update exposing (..)

import LiveView.Messages exposing (Msg(..))
import LiveView.Models as Model exposing (..)


update : Msg -> Model -> ( Model, Cmd Msg )
update message model =
    case message of
        NoOp ->
            ( model, Cmd.none )

        FileChanged fileChangedEvent ->
            let
                _ =
                    Debug.log "LiveView  received" fileChangedEvent

                oldState =
                    case model.state of
                        Initial ->
                            { serverMessages = [fileChangedEvent]
                            }
                        ReceivingEvents s ->
                            {
                                serverMessages = fileChangedEvent :: s.serverMessages
                            }
                
            in
                ( {model | state = ReceivingEvents oldState}, Cmd.none )
