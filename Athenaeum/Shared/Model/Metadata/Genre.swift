/**
 Genre.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

enum Genre: String, Codable {
    case sciFi = "Science Fiction"
    case fantasy = "Fantasy"
    case ya = "Yound Adult"
    case romance = "Romance"
    case thriller = "Thriller"
    case mystery = "Mystery"
    case lit = "Literary"
    case childrens = "Children's"
    case nonFiction = "Non-fiction"
}
