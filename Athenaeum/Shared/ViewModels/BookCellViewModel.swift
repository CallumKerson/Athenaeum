/**
 BookCellViewModel.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Combine
import Foundation

class BookCellViewModel: ObservableObject, Identifiable {
    var subscriptions: Set<AnyCancellable> = []
    let id: String

    @Published var loading: Bool = true
    @Published var error: AthenaeumError?
    @Published var book: Book?

    init(id: String) {
        self.id = id
    }

    func reload() {
        Athenaeum
            .loadItem(withId: self.id)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] value in
                guard let self = self else { return }
                if case let .failure(error) = value {
                    self.error = error
                }
                self.loading = false
            }, receiveValue: { [weak self] item in
                guard let self = self else { return }
                self.book = item
            })
            .store(in: &self.subscriptions)
    }
}
