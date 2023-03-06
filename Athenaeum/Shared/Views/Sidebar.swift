/**
 Sidebar.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

//
//  Sidebar.swift
//  Athenaeum (macOS)
//
//  Created by Callum Kerson on 21/04/2021.
//
import SwiftUI

struct Sidebar: View {
    @Binding var selectedBookId: String?
    var booksLogicController: BooksLogicController

    enum NavigationItem {
        case all
        case fantasy
        case sciFi
        case ya
    }

    @State private var selection: NavigationItem? = .all

    var body: some View {
        List(selection: $selection) {
            NavigationLink(destination: BooksListView(
                viewModel: BooksListViewModel(
                    booksLogicController: booksLogicController,
                    genre: nil
                ),
                selectedBookId: $selectedBookId
            )) {
                Label("All", systemImage: "books.vertical")
            }
            Divider()
            NavigationLink(destination: BooksListView(
                viewModel: BooksListViewModel(
                    booksLogicController: booksLogicController,
                    genre: .fantasy
                ),
                selectedBookId: $selectedBookId
            )) {
                Label(Genre.fantasy.rawValue, systemImage: Genre.fantasy.icon)
            }
            NavigationLink(destination: BooksListView(
                viewModel: BooksListViewModel(
                    booksLogicController: booksLogicController,
                    genre: .sciFi
                ),
                selectedBookId: $selectedBookId
            )) {
                Label(Genre.sciFi.rawValue, icon: {
                    Circle()
                })
            }
//            NavigationLink(destination: BooksListView(
//                viewModel: BooksListViewModel(booksLogicController: booksLogicController, genre: .ya),
//                selectedBookId: $selectedBookId
//            )) {
//                Label(Genre.ya.rawValue, icon:
//                        Circle()
            ////                               .fill(person.profileColor)
//                               .frame(width: 44, height: 44, alignment: .center)
//                               .overlay(Text("YA"))
//                      )
//            }
        }
        .listStyle(SidebarListStyle())
        .frame(minWidth: 150, idealWidth: 150, maxWidth: 200, maxHeight: .infinity)
        .padding(.top, 16)
    }
}
