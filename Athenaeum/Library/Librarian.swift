/**
 Import.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation
import GoodReadsKit

class Librarian<R, S> where R: Repository, R.EntityObject == AudiobookFile, S: PreferencesStore {
    let preferences: S
    let repository: R

    init(withPreferences preferences: S, withRepository repository: R) {
        self.preferences = preferences
        self.repository = repository
    }

    func importAudiobook(fileURL: URL) {
        log.info("Importing audiobook file from \(fileURL.path)")
        let newBook = AudiobookFile(fromFile: fileURL)
        var destination = self.preferences.libraryPath
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

        if fileURL.isSameIgnoringSandbox(as: destination) {
            log.debug("New Audiobook \(newBook) is already in the right place")
        } else {
            log.debug("Moving audiobook file to \(destination.path)")
            do {
                try FileManager.default
                    .moveItemCreatingIntermediaryDirectories(at: fileURL, to: destination)
                newBook.location = destination
            } catch {
                log.error("Cannot move \(fileURL.path) to \(destination.path)")
                log.error(error)
                return
            }
        }
        log.info("Adding audiobook \(newBook) to library")
        do {
            try self.repository.insert(item: newBook)
        } catch {
            log.error(error)
        }
    }

    func setUpLibraryPath() {
        do {
            try FileManager.default
                .createDirectoryIfDoesNotExist(atURL: self.preferences.libraryPath)
        } catch {
            log.error("Could not create library at path \(self.preferences.libraryPath.path)")
        }

        if self.preferences.useImportDirectory {
            log.info("Enabling auto import")
            DispatchQueue.global(qos: .userInitiated).async {
                self.importFilesInLibraryPath()
            }
        }
    }

    func importFilesInLibraryPath() {
        let existingAudiobooks = self.repository.items
            .map { $0.location.deSandboxedPath }

        var files = [URL]()
        if let enumerator = FileManager.default
            .enumerator(at: self.preferences.libraryPath,
                        includingPropertiesForKeys: [.isRegularFileKey],
                        options: [.skipsHiddenFiles,
                                  .skipsPackageDescendants]) {
            for case let fileURL as URL in enumerator {
                do {
                    let fileAttributes = try fileURL.resourceValues(forKeys: [.isRegularFileKey])
                    if fileAttributes.isRegularFile! {
                        if fileURL.pathExtension == "m4b" {
                            if !existingAudiobooks.contains(fileURL.deSandboxedPath) {
                                files.append(fileURL)
                            }
                        }
                    }
                } catch {
                    log.error(error)
                }
            }
        }
        for file in files {
            self.importAudiobook(fileURL: file)
        }
    }
}
