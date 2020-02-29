/**
 ContentView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct ContentView: View {
    @ObservedObject var viewModel: ContentViewModel

    init(_ viewModel: ContentViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        NavigationView {
            NavigationMasterView(NavigationMasterViewModel(store: store))

            if viewModel.selectedAudiobook != nil {
                NavigationDetailView(NavigationDetailViewModel(id: viewModel.selectedAudiobook!.id,
                                                               store: store))
            }
        }
        .frame(minWidth: 1100, minHeight: 500)
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView(ContentViewModel(store: sampleStore))
    }
}
