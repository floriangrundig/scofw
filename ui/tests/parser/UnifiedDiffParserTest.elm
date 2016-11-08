module Parser.UnifiedDiffParserTest exposing (..)

import Test exposing (..)
import Expect
import Fuzz exposing (..)
import String
import Parser.UnifiedDiffParser as UnifiedDiffParser
import CommonModels exposing (..)


all : Test
all =
    Test.concat
        [ parseSingleLogEntry
        , parseSeveralLogEntries
        ]



-- TODO test hunk with only one line change!


parseSingleLogEntry : Test
parseSingleLogEntry =
    test "Parse entry" <|
        \() ->
            let
                logEntry =
                    """2016/08/28 17:55:40 diff --git a/Makefile b/Makefile
index 8847fa04b4360aee350728146dba874fe440c055..0d6411a3c4a58249ea5ce2d9a170368475aaba34 100644
--- a/Makefile
+++ b/Makefile
@@ -13,3 +13,3 @@ build-go: build-go-linux build-go-darwin
 build-go-linux:
-   foo
+   bar

@@ -18 +18 @@ build-go2: build-go2-linux build-go2-darwin
-   foo2
+   bar2
"""

                expectedFileChanges =
                    [ Ok
                        { oldName = Ok "Makefile"
                        , newName = Ok "Makefile"
                        , op = Modified
                        , time = Ok "2016/08/28 17:55:40"
                        , hunks =
                            [ Ok
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
                            , Ok
                                { context = Just "build-go2: build-go2-linux build-go2-darwin"
                                , lines =
                                    [ Deletion "   foo2"
                                    , Addition "   bar2"
                                    , Context ""
                                    ]
                                , ranges =
                                    Ok
                                        { fromFileLineNumberStart = 18
                                        , toFileLineNumberStart = 18
                                        , fromFileLineNumberEnd = 18
                                        , toFileLineNumberEnd = 18
                                        }
                                , additions = 1
                                , deletions = 1
                                }
                            ]
                        }
                    ]
            in
                Expect.equal expectedFileChanges (UnifiedDiffParser.parse logEntry)


parseSeveralLogEntries : Test
parseSeveralLogEntries =
    test "Parse entry" <|
        \() ->
            let
                logEntry =
                    """2016/08/28 17:55:40 diff --git a/Makefile b/Makefile
index 8847fa04b4360aee350728146dba874fe440c055..0d6411a3c4a58249ea5ce2d9a170368475aaba34 100644
--- a/Makefile
+++ b/Makefile
@@ -13,3 +13,3 @@ build-go: build-go-linux build-go-darwin
 build-go-linux:
-   foo
+   bar

2016/08/28 18:55:40 diff --git a/Makefile2 b/Makefile2
index 8847fa04b4360aee350728146dba874fe440c055..0d6411a3c4a58249ea5ce2d9a170368475aaba34 100644
--- a/Makefile2
+++ b/Makefile2
@@ -13,3 +13,3 @@ build-go2: build-go-linux build-go-darwin
 build-go-linux2:
-   foo2
+   bar2

"""

                expectedFileChanges =
                    [ Ok
                        { oldName = Ok "Makefile"
                        , newName = Ok "Makefile"
                        , op = Modified
                        , time = Ok "2016/08/28 17:55:40"
                        , hunks =
                            [ Ok
                                { context = Just "build-go: build-go-linux build-go-darwin"
                                , lines =
                                    [ Context "build-go-linux:"
                                    , Deletion "   foo"
                                    , Addition "   bar"
                                    , Context ""
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
                            ]
                        }
                    , Ok
                        { oldName = Ok "Makefile2"
                        , newName = Ok "Makefile2"
                        , op = Modified
                        , time = Ok "2016/08/28 18:55:40"
                        , hunks =
                            [ Ok
                                { context = Just "build-go2: build-go-linux build-go-darwin"
                                , lines =
                                    [ Context "build-go-linux2:"
                                    , Deletion "   foo2"
                                    , Addition "   bar2"
                                    , Context ""
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
                            ]
                        }
                    ]
            in
                Expect.equal expectedFileChanges (UnifiedDiffParser.parse logEntry)
