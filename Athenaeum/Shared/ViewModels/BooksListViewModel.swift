/**
 BooksListViewModel.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Combine
import Foundation

public class BooksListViewModel: ObservableObject {
    var subscriptions: Set<AnyCancellable> = []
    private let genre: Genre

    @Published var loading: Bool = true
    @Published var error: AthenaeumError?
    @Published var items: [BookCellViewModel] = []

    init(genre: Genre) {
        self.genre = genre
    }

    func reload() {
        Athenaeum
            .loadItems()
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] value in
                guard let self = self else { return }
                if case let .failure(error) = value {
                    self.error = error
                }
                self.loading = false
            }, receiveValue: { [weak self] items in
                guard let self = self else { return }
                self.items = items
            })
            .store(in: &self.subscriptions)
    }
}
