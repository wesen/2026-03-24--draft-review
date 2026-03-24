import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

type View = "dashboard" | "articles" | "article-edit" | "article-settings" | "article-review";

interface UiState {
  view: View;
  selectedArticleId: string | null;
  activeModal: string | null;
  activeMenu: string | null;
}

const initialState: UiState = {
  view: "dashboard",
  selectedArticleId: null,
  activeModal: null,
  activeMenu: null,
};

export const uiSlice = createSlice({
  name: "ui",
  initialState,
  reducers: {
    setView(state, action: PayloadAction<View>) {
      state.view = action.payload;
    },
    selectArticle(state, action: PayloadAction<string>) {
      state.selectedArticleId = action.payload;
    },
    openModal(state, action: PayloadAction<string>) {
      state.activeModal = action.payload;
    },
    closeModal(state) {
      state.activeModal = null;
    },
    setActiveMenu(state, action: PayloadAction<string | null>) {
      state.activeMenu = action.payload;
    },
  },
});

export const { setView, selectArticle, openModal, closeModal, setActiveMenu } =
  uiSlice.actions;
