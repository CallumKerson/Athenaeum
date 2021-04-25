/**
 Athenaeum.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Alamofire
import Combine
import Foundation

final class Athenaeum {
    private static let headers: HTTPHeaders =
        [.authorization(username: "librarian", password: "a9ba69d1-fa46-485c-b410-e4055467210f")]
    private static let apiHost = "http://localhost:8030/api/v1"

    static func loadItems() -> AnyPublisher<[BookCellViewModel], AthenaeumError> {
        AF.request("\(self.apiHost)/books/", headers: self.headers)
            .publishDecodable(type: Response.self)
            .value()
            .map { response in
                response.books.map { book in
                    BookCellViewModel(id: book.id)
                }
            }
            .mapError { AthenaeumError.map($0) }
            .eraseToAnyPublisher()
    }

    static func loadItem(withId id: String) -> AnyPublisher<Book, AthenaeumError> {
        AF.request("\(self.apiHost)/books/\(id)", headers: self.headers)
            .publishDecodable(type: Book.self)
            .value()
            .mapError { AthenaeumError.map($0) }
            .eraseToAnyPublisher()
    }
}

struct Response: Codable {
    var books: [Book]
}
