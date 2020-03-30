/**
 MetadataService.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

struct MetadataService {
    private static let persistenceQueue =
        DispatchQueue(label: "com.umbra.Athenaeum.metadataQueue")

    private var didStateChangeCancellable: AnyCancellable?

    init(store: Store<GlobalAppState>) {
        self.didStateChangeCancellable = store.stateSubject.sink(receiveValue: { state in
            MetadataService.persistenceQueue.async {
                let recievedAudiobooks: [AudioBook] = Array(state.audiobookState.audiobooks.values)
                    .loadedAudiobooks
                for book in recievedAudiobooks {
                    if let metadata = book.metadata {
                        var metadataLocation = book.location
                        metadataLocation.deletePathExtension()
                        metadataLocation.appendPathExtension("json")

                        let encoder = JSONEncoder()
                        encoder.outputFormatting = [.prettyPrinted, .sortedKeys]
                        if let encodedData = try? encoder.encode(metadata) {
                            do {
                                try encodedData.write(to: metadataLocation)
                            } catch {
                                log
                                    .error("Failed to write JSON data: \(error.localizedDescription)")
                            }
                        }
                    }
                }
            }
        })
    }
}
