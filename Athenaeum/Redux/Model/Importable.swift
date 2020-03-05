/**
 Importable.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

enum Loadable<T: Equatable>: Equatable, Codable where T: Codable {
    case loading(T)
    case loaded(T)
    case errored(T, message: String)

    private enum CodingKeys: String, CodingKey {
        case loading
        case loaded
        case erroredAudiobook
        case errorMessage
    }

    enum LoadableCodingError: Error {
        case decoding(String)
    }

    init(from decoder: Decoder) throws {
        let values = try decoder.container(keyedBy: CodingKeys.self)
        if let value = try? values.decode(T.self, forKey: .loading) {
            self = .loading(value)
            return
        }
        if let value = try? values.decode(T.self, forKey: .loaded) {
            self = .loaded(value)
            return
        }
        if let value = try? values.decode(T.self, forKey: .erroredAudiobook),
            let error = try? values.decode(String.self, forKey: .errorMessage) {
            self = .errored(value, message: error)
            return
        }
        throw LoadableCodingError.decoding("Cannot decode \(dump(values))")
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        switch self {
        case let .loaded(audiobook):
            try container.encode(audiobook, forKey: .loaded)
        case let .loading(audiobook):
            try container.encode(audiobook, forKey: .loading)
        case let .errored(audiobook, message):
            try container.encode(audiobook, forKey: .erroredAudiobook)
            try container.encode(message, forKey: .errorMessage)
        }
    }
}

extension Loadable {
    var isLoaded: Bool {
        if case .loaded = self {
            return true
        }
        return false
    }

    var isLoading: Bool {
        if case .loading = self {
            return true
        }
        return false
    }
}

extension Array where Element == Loadable<AudioBook> {
    var loadedAudiobooks: [AudioBook] {
        self.compactMap { (loadable) -> AudioBook? in
            if case let .loaded(incomingAudiobook) = loadable {
                return incomingAudiobook
            }
            return nil
        }
    }
}