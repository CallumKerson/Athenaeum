/**
 Genre.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Foundation

enum Genre: String, CaseIterable, Identifiable {
    var id: String {
        rawValue
    }

    case all

    var name: String {
        switch self {
        case .all: return "All"
        }
    }

    var urlSuffix: String {
        self.rawValue
    }

    var icon: String {
        switch self {
        case .all: return "book"
        }
    }
}
