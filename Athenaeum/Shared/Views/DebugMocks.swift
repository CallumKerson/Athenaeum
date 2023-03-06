/**
 DebugMocks.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Combine
import Foundation

// #if DEBUG

class MockController: BooksLogicController {
    var networkController: NetworkController

    static var book: Book {
        Book(id: "001", title: "A New Book",
             author: [Person(givenNames: "An", familyName: "Author")],
             summary: "This is a book",
             releaseDate: ISO8601DateFormatter().date(from: "2021-04-25T08:00:00Z"),
             series: nil)
    }

    init() {
        self.networkController = FoundationNetworkController()
    }

    func getBooks() -> AnyPublisher<Books, AthenaeumError> {
        Just(Books(books: [MockController.book]))
            .setFailureType(to: AthenaeumError.self)
            .eraseToAnyPublisher()
    }

    func getBook(withID _: String) -> AnyPublisher<Book, AthenaeumError> {
        Just(MockController.book)
            .setFailureType(to: AthenaeumError.self)
            .eraseToAnyPublisher()
    }

    func getBooks(genre _: Genre) -> AnyPublisher<Books, AthenaeumError> {
        Just(Books(books: [MockController.book]))
            .setFailureType(to: AthenaeumError.self)
            .eraseToAnyPublisher()
    }

    func updateBook(withID _: String,
                    updatedBook _: EditBook) -> AnyPublisher<Book, AthenaeumError>
    {
        Just(MockController.book)
            .setFailureType(to: AthenaeumError.self)
            .eraseToAnyPublisher()
    }
}

// #endif
