/**
 BookTest.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

@testable import Athenaeum
import XCTest

class BookTest: XCTestCase {
    func testDecodeJson() throws {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601

        let jsonData = self.simpleGenreResult.data(using: .utf8)!
        let books: Books = try decoder.decode(Books.self, from: jsonData)
        XCTAssertEqual(books.books[0].title, "The Tyrant Baru Cormorant")
        XCTAssertEqual(books.books[0].genre?[0], .fantasy)
    }

    let simpleGenreResult = """
    {
        "books": [
            {
                "id": "001",
                "title": "The Tyrant Baru Cormorant",
                "author": [
                    {
                        "givenNames": "Seth",
                        "familyName": "Dickinson"
                    }
                ],
                "releaseDate": "2020-08-11T08:00:00Z",
                "genre": [
                    "Fantasy"
                ],
                "series": {
                    "entry": 3,
                    "title": "The Masquerade"
                }
            }
        ]
    }
    """
}
