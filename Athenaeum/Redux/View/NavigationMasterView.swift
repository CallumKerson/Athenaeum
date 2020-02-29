/**
 NavigationMasterView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct NavigationMasterView: View {
    @ObservedObject var viewModel: NavigationMasterViewModel

    init(_ viewModel: NavigationMasterViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        VStack {
            NavigationHeaderView(NavigationHeaderViewModel(store: self.viewModel.store))
                .padding([.top, .leading], 8)
                .padding(.trailing, 4)

            AudiobookListView(AudiobookListViewModel(store: self.viewModel.store))
                .listStyle(SidebarListStyle())
        }
        .frame(minWidth: 300, maxWidth: 300)
    }
}

#if DEBUG
    struct NavigationMasterView_Previews: PreviewProvider {
        static var previews: some View {
            NavigationMasterView(NavigationMasterViewModel(store: sampleStore))
        }
    }
#endif