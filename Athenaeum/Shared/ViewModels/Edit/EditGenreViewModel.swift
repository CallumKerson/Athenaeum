/**
 EditGenreViewModel.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Foundation

class EditGenreViewModel: ObservableObject {
    @Published var genreState: GenreState

    init(initialGenre: Genre?) {
        if let initialGenre = initialGenre {
            switch initialGenre {
            case .fantasy:
                self.genreState = .fantasy
            case .sciFi:
                self.genreState = .sciFi
            }
        } else {
            self.genreState = .none
        }
    }

    func genre() -> [Genre]? {
        switch self.genreState {
        case .none:
            return nil
        case .fantasy:
            return [.fantasy]
        case .sciFi:
            return [.sciFi]
        }
    }
}

enum GenreState: String, CaseIterable, Identifiable, Hashable, CustomStringConvertible {
    var id: String {
        rawValue
    }

    var description: String {
        rawValue
    }

    case none = "-- Select a genre --"
    case sciFi = "Science Fiction"
    case fantasy = "Fantasy"
}
