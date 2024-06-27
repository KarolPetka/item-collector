import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {HttpClientModule} from '@angular/common/http';
import {Routing} from './routing';

import {AppComponent} from './app.component';
import {MainViewComponent} from './pages/main-view/main-view.component';
import {NewCollectionComponent} from './pages/new-collection/new-collection.component';
import {EditCollectionComponent} from './pages/edit-collection/edit-collection.component';
import {NewItemComponent} from './pages/new-item/new-item.component';
import {EditItemComponent} from './pages/edit-item/edit-item.component';
import {LoginComponent} from "./pages/login/login.component";
import {FormsModule} from "@angular/forms";
import {SignupComponent} from "./pages/signup/signup.component";

@NgModule({
  declarations: [
    AppComponent,
    SignupComponent,
    LoginComponent,
    MainViewComponent,
    NewCollectionComponent,
    EditCollectionComponent,
    NewItemComponent,
    EditItemComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    Routing,
    FormsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {

}
