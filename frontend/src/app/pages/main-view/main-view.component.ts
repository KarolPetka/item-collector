import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Params, Router} from '@angular/router';
import {Service} from 'src/app/service';

@Component({
  selector: 'app-item-view',
  templateUrl: './main-view.component.html',
  styleUrls: ['./main-view.component.css']
})
export class MainViewComponent implements OnInit {
  collections: any;
  items: any;

  selectedCollectionId: any;
  token: any;

  constructor(private service: Service, private route: ActivatedRoute, private router: Router) {
  }

  ngOnInit() {
    this.token = localStorage.getItem('jwtToken') as string;

    if (this.token === null) {
      this.router.navigate(['/login']);
    }

    this.route.params.subscribe(
      (params: Params) => {
        if (params['collectionId']) {
          this.selectedCollectionId = params['collectionId'];
          this.service.getItems(params['collectionId'], this.token).subscribe((items: any) => {
            this.items = items;
          })
        } else {
          this.items = undefined;
        }
      }
    )

    this.service.getCollections(this.token).subscribe((collections: any) => {
      this.collections = collections;
    })
  }

  onItemClick(item: any) {
    item.isCollected = !item.isCollected;
    this.service.updateItem(item, item.collectionId, item.id, this.token).subscribe(() => {
    })
  }

  onDeleteCollectionClick() {
    this.service.deleteCollection(this.selectedCollectionId, this.token).subscribe((res: any) => {
      this.router.navigate(['/collections']);
    })
  }

  onDeleteItemClick(id: string) {
    this.service.deleteItem(this.selectedCollectionId, id, this.token).subscribe((res: any) => {
      this.items = this.items.filter((val: { id: string; }) => val.id !== id);
    })
  }

  isCollectionIdNotSet(): boolean {
    return this.selectedCollectionId === undefined;
  }

  getRarityImage(rarity: number): string {
    if (rarity === 1) {
      return '../../../assets/high-rarity.svg';
    } else if (rarity === 2) {
      return '../../../assets/medium-rarity.svg';
    } else {
      return '../../../assets/low-rarity.svg';
    }
  }

  logOut(){
    localStorage.removeItem('jwtToken');
    this.router.navigate(['/login']);
  }
}
