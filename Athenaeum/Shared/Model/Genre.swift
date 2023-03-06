/**
 Genre.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Foundation

enum Genre: String, CaseIterable, Identifiable, Codable {
    var id: String {
        rawValue
    }

    case sciFi = "Science Fiction"
    case fantasy = "Fantasy"
//    case childrens = "Children's"
//    case ya = "Young Adult"
//    case nonfiction = "Non-Fiction"
//    case historical = "Historical"

    var urlSuffix: String {
        self.rawValue
    }

    var icon: String {
        switch self {
        case .sciFi: return "bolt"
        case .fantasy: return "wand.and.stars"
//        case .childrens: return "face.smiling"
//        case .ya: return ""
//        case .nonfiction: return ""
//        case .historical: return "clock"
        }
    }
}
