//
//  AddressBookViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/22.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "AddressBookViewController.h"

@interface AddressBookViewController () <UITableViewDelegate, UITableViewDataSource, AddressBookTopViewDelegate, FriendCollectDelegate>

@property (strong, nonatomic) AddressBookTopView *topView;

@property (strong, nonatomic) UITableView *tabel;

@property (strong, nonatomic) UITableView *listTabel;

@property (strong, nonatomic) NSArray *arrNewList, *arrOldList;
@property (strong, nonatomic) NSArray *listData;

@end

@implementation AddressBookViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor colorWithRed:240 / 255.0 green:240 / 255.0 blue:240 / 255.0 alpha:1];
    [self requestData];
    [self requestListData];
    [self creatTopView];
    [self.view addSubview:self.tabel];
    [self.view addSubview:self.listTabel];
}

- (void)requestData {
    [HttpData requestFriendDataToken:@"f2e0c88f-7d32-3464-9cc5" success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"]integerValue] == 200) {
                self.arrNewList = json[@"data"][@"new_user_list"];
                self.arrOldList = json[@"data"][@"old_user_list"];
                [self.tabel reloadData];
            }
        }
        else {
            
        }
    } failure:^(NSError * _Nonnull err) {
        
    }];
}

- (void)requestListData {
    [HttpData requestCollectFriendDataToken:@"f2e0c88f-7d32-3464-9cc5" success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"]integerValue] == 200) {
                self.listData = json[@"data"];
                [self.listTabel reloadData];
            }
        }
    } failure:^(NSError * _Nonnull err) {
        
    }];
}

- (void)creatTopView {
    _topView = [[AddressBookTopView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT / 4)];
    _topView.delegate = self;
    _topView.backgroundColor = [UIColor whiteColor];
    [self.view addSubview:self.topView];
}

- (UITableView*)tabel {
    if (!_tabel) {
        _tabel = [[UITableView alloc] initWithFrame:CGRectMake(0, self.topView.frame.size.height + 10, SCREEN_WIDTH, SCREEN_HEIGHT - self.topView.frame.size.height - 10) style:UITableViewStyleGrouped];
        _tabel.delegate = self;
        _tabel.dataSource = self;
        _tabel.backgroundColor = [UIColor colorWithRed:245 / 250.0 green:245 / 250.0 blue:245 / 250.0 alpha:1];
        _tabel.separatorStyle = UITableViewCellSelectionStyleNone;
        _tabel.tag = 1;
        _tabel.hidden = NO;
    }
    return _tabel;
}

- (UITableView*)listTabel {
    if (!_listTabel) {
        _listTabel = [[UITableView alloc] initWithFrame:CGRectMake(0, self.topView.frame.size.height + 10, SCREEN_WIDTH, SCREEN_HEIGHT - self.topView.frame.size.height - 10) style:UITableViewStyleGrouped];
        _listTabel.delegate = self;
        _listTabel.dataSource = self;
        _listTabel.backgroundColor = [UIColor colorWithRed:245 / 250.0 green:245 / 250.0 blue:245 / 250.0 alpha:1];
        _listTabel.separatorStyle = UITableViewCellSelectionStyleNone;
        _listTabel.tag = 2;
        _listTabel.hidden = YES;
    }
    return _listTabel;
}

- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView {
    return tableView.tag == 1 ? 2 : 1;
}

- (NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section {
    if (tableView.tag == 1) return section == 0 ? 1 : self.arrOldList.count;
    else return self.listData.count;
}

- (CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath {
    if (tableView.tag == 1) return indexPath.section == 0 ? SCREEN_WIDTH / 3 + 20 : 120;
    else return 100;
}

- (CGFloat)tableView:(UITableView *)tableView heightForHeaderInSection:(NSInteger)section {
    return tableView.tag == 1 ? 45 : 0.01;
}

- (CGFloat)tableView:(UITableView *)tableView heightForFooterInSection:(NSInteger)section {
    return 0.01;
}

- (UIView *)tableView:(UITableView *)tableView viewForHeaderInSection:(NSInteger)section {
    if (tableView.tag == 2) return nil;
    UIView *view = [[UIView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, 45)];
    view.backgroundColor = [UIColor whiteColor];
    UILabel *label = [[UILabel alloc] initWithFrame:CGRectMake(15, 0, SCREEN_WIDTH, 45)];
    label.text = section == 0 ? @"New Matches" : @"Messages";
    label.textColor = [UIColor colorWithRed:178 / 255.0 green:34 / 255.0 blue:34 / 255.0 alpha:1];
    label.font = BoldFont(20);
    [view addSubview:label];
    return view;
}

- (UIView *)tableView:(UITableView *)tableView viewForFooterInSection:(NSInteger)section {
    return nil;
}

- (UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath {
    if (tableView.tag == 1) {
        if (indexPath.section == 0) {
            static NSString *cellid = @"cellid";
            OneFriendTableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"cellid"];
            if (!cell) {
                cell = [[OneFriendTableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
            }
            cell.delegate = self;
            cell.collList = self.arrNewList;
            [cell.collection reloadData];
            
            cell.selectionStyle = UITableViewCellSelectionStyleNone;
            return cell;
        }
        else {
            static NSString *cellid = @"cellid";
            FriendTableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"cellid"];
            if (!cell) {
                cell = [[FriendTableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
            }
//            cell.headImg.image = [UIImage imageNamed:@"navLeftImg"];
//            cell.nameLab.text = @"Crayon Shinchan";
//            cell.messageLab.text = @"Nice to meet you!";
            [cell.headImg sd_setImageWithURL:self.arrOldList[indexPath.row][@"photo"] placeholderImage:[UIImage imageNamed:@"navLeftImg"]];
            cell.nameLab.text = self.arrOldList[indexPath.row][@"name"];
            cell.messageLab.text = self.arrOldList[indexPath.row][@"message"];
            
            cell.selectionStyle = UITableViewCellSelectionStyleNone;
            return cell;
        }
    }
    else {
        static NSString *cellid = @"cellid";
        FriendTableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"cellid"];
        if (!cell) {
            cell = [[FriendTableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
        }
//        cell.headImg.image = [UIImage imageNamed:@"navLeftImg"];
//        cell.nameLab.text = @"Lisa";
//        cell.messageLab.text = @"Nice to meet you!";
        
        [cell.headImg sd_setImageWithURL:_listData[indexPath.row][@"photo"] placeholderImage:[UIImage imageNamed:@"navLeftImg"]];
        cell.nameLab.text = _listData[indexPath.row][@"name"];
        cell.messageLab.text = _listData[indexPath.row][@"time"];
        
        cell.selectionStyle = UITableViewCellSelectionStyleNone;
        return cell;
    }
}

- (void)clickAddressBookTopViewButton:(NSInteger)tag {
    switch (tag) {
        case 1:
            [self.navigationController popViewControllerAnimated:YES];
            break;
        case 2:
            [self messagesBtnClick];
            break;
        case 3:
            [self listBtnClick];
        default:
            break;
    }
}

- (void)messagesBtnClick {
    _tabel.hidden = NO;
    _listTabel.hidden = YES;
}

- (void)listBtnClick {
    _tabel.hidden = YES;
    _listTabel.hidden = NO;
}

- (void)clickFriendCollectiomItemTag:(NSInteger)tag {
//    NSLog(@"%ld", tag);
    ChatViewController *chatVC = [[ChatViewController alloc] init];
    chatVC->url = self.arrNewList[tag][@"photo"];
    chatVC->name = self.arrNewList[tag][@"name"];
    [self.navigationController pushViewController:chatVC animated:YES];
    
}

- (void)tableView:(UITableView *)tableView didSelectRowAtIndexPath:(NSIndexPath *)indexPath {
//    NSLog(@"%ld", tableView.tag);
//    NSLog(@"%ld", indexPath.section);
//    NSLog(@"%ld", indexPath.row);
    ChatViewController *chatVC = [[ChatViewController alloc] init];
    [self.navigationController pushViewController:chatVC animated:YES];
    chatVC->url = self.arrOldList[indexPath.row][@"photo"];
    chatVC->name = self.arrOldList[indexPath.row][@"name"];
    
}


@end
