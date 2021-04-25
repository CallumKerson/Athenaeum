/**
 AthenaeumApp.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import SwiftUI

@main
struct AthenaeumApp: App {
    @State private var selectedBookId: String?

    var body: some Scene {
        WindowGroup {
            #if os(macOS)
                NavigationView {
                    Sidebar(selectedBookId: $selectedBookId)
                    BooksListView(
                        viewModel: BooksListViewModel(genre: .all),
                        selectedBookId: $selectedBookId
                    )
                    Text("Select book...")
                }
            #else
                NavigationView {
                    TabBar()
                        .navigationTitle("Books")
                }.navigationViewStyle(StackNavigationViewStyle())
            #endif
        }
    }
}
