/**
 FileLibrarian.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

class FileLibrarian: Librarian<RealmRepository<AudiobookFile>, UserDefaultsPreferencesStore> {
    static var global =
        FileLibrarian(withPreferences: UserDefaultsPreferencesStore.global,
                      withRepository: AudiobookRepository.global)
}
