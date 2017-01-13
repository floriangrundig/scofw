module Tests exposing (..)

import Test exposing (..)
import Parser.UnifiedDiffParserTest as UnifiedDiffParserTest
import LiveView.ViewTest as LiveViewTest

all : Test
all =
    Test.concat
        [ UnifiedDiffParserTest.all
        , LiveViewTest.all
        ]
