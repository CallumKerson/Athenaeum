/**
 BooksLogicController.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Combine
import Foundation

protocol BooksLogicController: AnyObject {
    var networkController: NetworkController { get }

    func getBooks() -> AnyPublisher<Books, AthenaeumError>
    func getBook(withID id: String) -> AnyPublisher<Book, AthenaeumError>
    func getBooks(genre: Genre) -> AnyPublisher<Books, AthenaeumError>
    func updateBook(withID id: String, updatedBook: EditBook) -> AnyPublisher<Book, AthenaeumError>
}

final class BooksEndpointLogicController: BooksLogicController {
    let networkController: NetworkController

    init(networkController: NetworkController) {
        self.networkController = networkController
    }

    func getBooks() -> AnyPublisher<Books, AthenaeumError> {
        let endpoint = Endpoint.books

        return self.networkController.get(type: Books.self,
                                          url: endpoint.athenaeumURL,
                                          headers: endpoint.headers)
    }

    func getBook(withID id: String) -> AnyPublisher<Book, AthenaeumError> {
        let endpoint = Endpoint.book(withID: id)

        return self.networkController.get(type: Book.self,
                                          url: endpoint.athenaeumURL,
                                          headers: endpoint.headers)
    }

    func getBooks(genre: Genre) -> AnyPublisher<Books, AthenaeumError> {
        let endpoint = Endpoint.genre(genre)

        return self.networkController.get(type: Books.self,
                                          url: endpoint.athenaeumURL,
                                          headers: endpoint.headers)
    }

    func updateBook(withID id: String,
                    updatedBook: EditBook) -> AnyPublisher<Book, AthenaeumError>
    {
        let endpoint = Endpoint.book(withID: id)
        return self.networkController.patch(updateObject: updatedBook,
                                            returnType: Book.self,
                                            url: endpoint.athenaeumURL,
                                            headers: endpoint.headers)
    }
}
