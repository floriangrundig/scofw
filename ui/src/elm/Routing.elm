module Routing exposing (..)

import Navigation
import UrlParser exposing (Parser, (</>), s, int, string, map, oneOf, parseHash)


type Route
    = RootRoute
    | LiveViewRoute
    | NotFoundRoute


routes : Parser (Route -> a) a
routes =
    oneOf
        [ map RootRoute (s "")
        , map LiveViewRoute (s "live")
        ]


parseHash : Navigation.Location -> Route
parseHash location =
    case ( location.pathname, location.search, location.hash ) of
        ( "/", "", "" ) ->
            RootRoute

        _ ->
            UrlParser.parseHash routes location |> Maybe.withDefault NotFoundRoute


parse : Navigation.Location -> Route
parse location =
    UrlParser.parsePath routes location |> Maybe.withDefault NotFoundRoute
