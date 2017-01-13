module LiveView.ViewTest exposing (..)

import Test exposing (..)
import Expect
import CommonModels exposing (..)
import LiveView.View as View exposing (..)


all : Test
all =
    Test.concat
        [ transformHunkToViewableHunk1
        , transformHunkToViewableHunk2
        , transformHunkToViewableHunk3
        , transformHunkToViewableHunk4
        ]


transformHunkToViewableHunk1 : Test
transformHunkToViewableHunk1 =
    test "transform hunk to viewable hunk lines - 1" <|
        \() ->
            let
                hunk =
                    { context = Just "build-go: build-go-linux build-go-darwin"
                    , lines =
                        [ Context "build-go-linux:"
                        , Deletion "   foo"
                        , Addition "   bar"
                        , Context ""
                        ]
                    , ranges =
                        Ok
                            { fromFileLineNumberStart = 13
                            , toFileLineNumberStart = 13
                            , fromFileLineNumberEnd = 3
                            , toFileLineNumberEnd = 3
                            }
                    , additions = 1
                    , deletions = 1
                    }
            in
                Expect.equal
                    ( [ Context_ ( 13, "build-go-linux:" )
                      , Addition_ ( 14, "   bar" )
                      , Context_ ( 15, "" )
                      ]
                    , [ Context_ ( 13, "build-go-linux:" )
                      , Deletion_ ( 14, "   foo" )
                      , Context_ ( 15, "" )
                      ]
                    )
                    (View.asViewableHunkLines hunk)


transformHunkToViewableHunk2 : Test
transformHunkToViewableHunk2 =
    test "transform hunk to viewable hunk lines - 2" <|
        \() ->
            let
                hunk =
                    { context = Just "build-go: build-go-linux build-go-darwin"
                    , lines =
                        [ Context "build-go-linux:"
                        , Deletion "   foo1"
                        , Deletion "   foo2"
                        , Addition "   bar1"
                        , Context ""
                        ]
                    , ranges =
                        Ok
                            { fromFileLineNumberStart = 13
                            , toFileLineNumberStart = 13
                            , fromFileLineNumberEnd = 3
                            , toFileLineNumberEnd = 3
                            }
                    , additions = 1
                    , deletions = 1
                    }
            in
                Expect.equal
                    ( [ Context_ ( 13, "build-go-linux:" )
                      , Addition_ ( 14, "   bar1" )
                      , AdditionPadding
                      , Context_ ( 15, "" )
                      ]
                    , [ Context_ ( 13, "build-go-linux:" )
                      , Deletion_ ( 14, "   foo1" )
                      , Deletion_ ( 15, "   foo2" )
                      , Context_ ( 16, "" )
                      ]
                    )
                    (View.asViewableHunkLines hunk)

transformHunkToViewableHunk3 : Test
transformHunkToViewableHunk3 =
    test "transform hunk to viewable hunk lines - 3" <|
        \() ->
            let
                hunk =
                    { context = Just "build-go: build-go-linux build-go-darwin"
                    , lines =
                        [ Context "build-go-linux:"
                        , Context " C1"
                        , Deletion "   foo1"
                        , Deletion "   foo2"
                        , Addition "   bar1"
                        , Addition "   bar2"
                        , Addition "   bar3"
                        , Context " C2"
                        ]
                    , ranges =
                        Ok
                            { fromFileLineNumberStart = 13
                            , toFileLineNumberStart = 13
                            , fromFileLineNumberEnd = 3
                            , toFileLineNumberEnd = 3
                            }
                    , additions = 1
                    , deletions = 1
                    }
            in
                Expect.equal
                    ( [ Context_ ( 13, "build-go-linux:" )
                      , Context_ ( 14, " C1" )
                      , Addition_ ( 15, "   bar1" )
                      , Addition_ ( 16, "   bar2" )
                      , Addition_ ( 17, "   bar3" )
                      , Context_ ( 18, " C2" )
                      ]
                    , [ Context_ ( 13, "build-go-linux:" )
                      , Context_ ( 14, " C1" )
                      , Deletion_ ( 15, "   foo1" )
                      , Deletion_ ( 16, "   foo2" )
                      , DeletionPadding
                      , Context_ ( 17, " C2" )
                      ]
                    )
                    (View.asViewableHunkLines hunk)


transformHunkToViewableHunk4 : Test
transformHunkToViewableHunk4 =
    test "transform hunk to viewable hunk lines - 4" <|
        \() ->
            let
                hunk =
                    { context = Just "build-go: build-go-linux build-go-darwin"
                    , lines =
                        [ Context "build-go-linux:"
                        , Context " C1"
                        , Addition "   bar1"
                        , Addition "   bar2"
                        , Addition "   bar3"
                        , Context " C2"
                        , Deletion "   foo1"
                        , Deletion "   foo2"
                        ]
                    , ranges =
                        Ok
                            { fromFileLineNumberStart = 13
                            , toFileLineNumberStart = 13
                            , fromFileLineNumberEnd = 3
                            , toFileLineNumberEnd = 3
                            }
                    , additions = 1
                    , deletions = 1
                    }
            in
                Expect.equal
                    ( [ Context_ ( 13, "build-go-linux:" )
                      , Context_ ( 14, " C1" )
                      , Addition_ ( 15, "   bar1" )
                      , Addition_ ( 16, "   bar2" )
                      , Addition_ ( 17, "   bar3" )
                      , Context_ ( 18, " C2" )
                      , AdditionPadding
                      , AdditionPadding
                      ]
                    , [ Context_ ( 13, "build-go-linux:" )
                      , Context_ ( 14, " C1" )
                      , DeletionPadding
                      , DeletionPadding
                      , DeletionPadding
                      , Context_ ( 15, " C2" )
                      , Deletion_ ( 16, "   foo1" )
                      , Deletion_ ( 17, "   foo2" )
                      ]
                    )
                    (View.asViewableHunkLines hunk)
