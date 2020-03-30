/**
 PersistenceService.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

struct PersistenceService {
    private static let persistenceQueue =
        DispatchQueue(label: "com.umbra.Athenaeum.persistenceQueue")

    private var didStateChangeCancellable: AnyCancellable?

    /// Initialised a persistence store that subscribes to app state changes from the store and writes them to a
    /// userdata.json file to save them.
    /// - Parameter store: Store that publishes app state changes
    init(store: Store<GlobalAppState>) {
        self.didStateChangeCancellable = store.stateSubject.sink(receiveValue: { state in
            PersistenceService.persistenceQueue.async {
                let encoder = JSONEncoder()
                encoder.outputFormatting = [.prettyPrinted, .sortedKeys]
                if let encodedData = try? encoder.encode(state) {
                    do {
                        try encodedData.write(to: PersistenceService.getSaveURL())
                    } catch {
                        log.error("Failed to write JSON data: \(error.localizedDescription)")
                    }
                }
            }
        })
    }

    /// Gets the URL for the user data persistence file
    static func getSaveURL() -> URL {
        do {
            var saveURL = try FileManager.default.url(for: .applicationSupportDirectory,
                                                      in: .userDomainMask,
                                                      appropriateFor: nil,
                                                      create: false)
            saveURL.appendPathComponent("Athenaeum", isDirectory: true)
            return saveURL.appendingPathComponent("userdata.json", isDirectory: false)

        } catch {
            fatalError("Cannot get save location \(error)")
        }
    }
}
