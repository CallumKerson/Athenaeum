/**
 AthenaeumApp.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Logging
import SwiftUI

let logger = Logger(label: "com.umbra.Athenaeum")

@main
struct AthenaeumApp: App {
    @State private var selectedBookId: String?
    var booksLogicController: BooksLogicController

    init() {
        logger.notice("Initialising Athenaeum")
        self
            .booksLogicController =
            BooksEndpointLogicController(networkController: FoundationNetworkController())
    }

    var body: some Scene {
        WindowGroup {
            //            #if os(macOS)
            NavigationView {
                Sidebar(selectedBookId: $selectedBookId, booksLogicController: booksLogicController)
                BooksListView(
                    viewModel: BooksListViewModel(
                        booksLogicController: booksLogicController,
                        genre: nil
                    ),
                    selectedBookId: $selectedBookId
                )
                Text("Select book...")
            }.frame(
                minWidth: 900,
                idealWidth: 1200,
                maxWidth: .infinity,
                minHeight: 600,
                idealHeight: 900,
                maxHeight: .infinity
            )
            //            #else
            //                NavigationView {
            //                    TabBar()
            //                        .navigationTitle("Books")
            //                }.navigationViewStyle(StackNavigationViewStyle())
            //            #endif
        }

        #if os(macOS)
            Settings {
                GeneralSettingsView()
            }
        #endif
    }
}
