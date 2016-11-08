module Tests exposing (..)

import Test exposing (..)
import Parser.UnifiedDiffParserTest as UnifiedDiffParserTest


all : Test
all =
    Test.concat
        [ UnifiedDiffParserTest.all
        ]
