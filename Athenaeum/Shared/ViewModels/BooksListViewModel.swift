/**
 BooksListViewModel.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Combine
import Foundation

public class BooksListViewModel: ObservableObject {
    var subscriptions: Set<AnyCancellable> = []
    private let genre: Genre?
    let booksLogicController: BooksLogicController

    @Published var loading: Bool = true
    @Published var error: AthenaeumError?
    @Published var books: Books = .init(books: [])

    init(booksLogicController: BooksLogicController, genre: Genre?) {
        self.booksLogicController = booksLogicController
        self.genre = genre
    }

    func reload() {
        if let genre = genre {
            self.booksLogicController
                .getBooks(genre: genre)
                .receive(on: DispatchQueue.main)
                .sink(receiveCompletion: { completion in
                    switch completion {
                    case let .failure(error):
                        logger
                            .error(
                                "Failed to reload book list for genre \(genre.rawValue): \(error)"
                            )
                    case .finished:
                        self.loading = false
                    }
                }) { books in
                    logger.info("Reloading List of books for genre \(genre.rawValue)")
                    self.books = books
                }
                .store(in: &self.subscriptions)
        } else {
            self.booksLogicController
                .getBooks()
                .receive(on: DispatchQueue.main)
                .sink(receiveCompletion: { completion in
                    switch completion {
                    case let .failure(error):
                        logger.error("Failed to reload book list: \(error)")
                        self.error = error
                        self.loading = false
                    case .finished:
                        self.loading = false
                    }
                }) { books in
                    logger.info("Reloading List of books")
                    self.books = books
                }
                .store(in: &self.subscriptions)
        }
    }
}
