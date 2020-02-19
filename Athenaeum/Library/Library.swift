/**
 Library.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

class Library {
    static var global = Library()

    let libraryURL: URL
    var audiobooks: [Audiobook] {
        getRepository().getAll()
    }

    init() {
        log.info("Initialising library")
        libraryURL = Preferences.libraryPath()
    }

    private func getRepository() -> RealmRepository<Audiobook> {
        RealmRepository()
    }
}
