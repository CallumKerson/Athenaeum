/**
 Import.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Cocoa
import Foundation

struct Import {
    let preferences: PreferencesStore
    let library: RepositoryLibrary

    init(withPreferences preferences: PreferencesStore = PreferencesStore.global,
         withLibrary library: RepositoryLibrary = RepositoryLibrary.global) {
        self.preferences = preferences
        self.library = library
    }

    func openImportAudiobookDialog() {
        let openPanel = NSOpenPanel()
        openPanel.canChooseFiles = true
        openPanel.allowsMultipleSelection = true
        openPanel.canChooseDirectories = false
        openPanel.canCreateDirectories = false
        openPanel.title = "Import Audiobooks"
        openPanel.allowedFileTypes = ["m4b"]

        openPanel.begin { response in
            if response == .OK {
                // TODO: Thread properly - this freezes the import window until import is complete
                for url in openPanel.urls {
                    DispatchQueue.global(qos: .userInitiated).async {
                        log.debug("User opening file \(url.path)")
                        self.importAudiobook(fileURL: url)
                    }
                }
            }
            openPanel.close()
        }
    }

    func importAudiobook(fileURL: URL) {
        log.info("Importing audiobook file from \(fileURL.path)")
        let newBook = AudiobookFile(fromFile: fileURL)
        var destination = preferences.libraryPath
            .appendingPathComponent(newBook.author, isDirectory: true)
        if let series = newBook.series {
            destination = destination
                .appendingPathComponent(series.title, isDirectory: true)
                .appendingPathComponent("\(series.entry) \(newBook.title)")
                .appendingPathExtension("m4b")
        } else {
            destination = destination
                .appendingPathComponent(newBook.title)
                .appendingPathExtension("m4b")
        }
        log.debug("Moving audiobook file to \(destination.path)")
        try! FileManager.default.moveItemCreatingIntermediaryDirectories(at: fileURL, to: destination)
        newBook.location = destination
        log.info("Adding audiobook \(newBook) to library")
        library.shelve(book: newBook)
    }
}
