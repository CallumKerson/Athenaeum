/**
 BookDetailViewModel.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Combine
import Foundation

class BookDetailViewModel: ObservableObject, Identifiable {
    var subscriptions: Set<AnyCancellable> = []
    let id: String
    let booksLogicController: BooksLogicController

    @Published var loading: Bool = true
    @Published var error: AthenaeumError?
    @Published var book: Book?

    init(booksLogicController: BooksLogicController, id: String) {
        self.id = id
        self.booksLogicController = booksLogicController
    }

    func reload() {
        self.booksLogicController
            .getBook(withID: self.id)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] value in
                guard let self = self else { return }
                if case let .failure(error) = value {
                    self.error = error
                    logger
                        .error(
                            "Failed to reload book details for book with id \(self.id): \(error)"
                        )
                }
                self.loading = false
            }, receiveValue: { [weak self] item in
                guard let self = self else { return }
                self.book = item
                logger.info("Reloading Book Details for book with id \(item.id)")
            })
            .store(in: &self.subscriptions)
    }
}
