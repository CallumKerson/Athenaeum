/**
 Library.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation
import Combine

final class Library: ObservableObject {
    
    static var global = Library()
    
    @Published var 🎧📚: [Audiobook] = []
    
    let repository = getRepository()
    
    private var didChangeCancellable: AnyCancellable?
    
    init() {
        log.info("Initialising Library")
        self.🎧📚 = self.repository.getAll()
        didChangeCancellable = repository.publisher.sink(receiveValue: { action in
            log.debug("Refreshing library due to reciving persistence action \(action.rawValue)")
            self.🎧📚 = self.repository.getAll()
        })
    }

    private static func getRepository() -> RealmRepository<Audiobook> {
        RealmRepository()
    }
}
