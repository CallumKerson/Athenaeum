/**
 Library.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation
import RealmSwift

class Library {
    static var global = Library()

    let libraryURL: URL
    var audiobooks: [Audiobook]

    var series: Set<Series>

    init() {
        log.info("Initialising library")
        libraryURL = Preferences.libraryPath()

        series = Set()
        audiobooks = []
    }

    func importAudiobook(fileURL: URL) {
        log.info("Importing audiobook file from \(fileURL.path)")
        let newBook = Audiobook.getBookFromFile(fileURL: fileURL)
        var destination = libraryURL
            .appendingPathComponent(newBook.author, isDirectory: true)
        if let seriesEntry = newBook.entry, let newBookSeries = newBook.series {
            destination = destination
                .appendingPathComponent(newBookSeries.title, isDirectory: true)
                .appendingPathComponent("\(seriesEntry) \(newBook.title)")
                .appendingPathExtension("m4b")
        } else {
            destination = destination
                .appendingPathComponent(newBook.title)
                .appendingPathExtension("m4b")
        }
        log.debug("Moving audiobook file to \(destination.path)")
        try! FileManager.default.moveItemCreatingIntermediaryDirectories(at: fileURL, to: destination)
        newBook.file = destination
    }

    private func getAudiobooks() -> [Audiobook] {
        var books = [Audiobook]()

        let enumerator = FileManager.default.enumerator(atPath: libraryURL.path)
        let filePaths = enumerator?.allObjects as! [String]
        let audiobookFilePaths = filePaths.filter { $0.contains(".m4b") }
        for audiobookPath in audiobookFilePaths {
            books.append(Audiobook.getBookFromFile(path: "\(libraryURL.path)/\(audiobookPath)"))
        }

        books.sort {
            $0.author < $1.author
        }

        return books
    }
}
