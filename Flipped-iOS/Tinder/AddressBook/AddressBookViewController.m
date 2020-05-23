//
//  AddressBookViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/22.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "AddressBookViewController.h"

@interface AddressBookViewController () <UITableViewDelegate, UITableViewDataSource>

@property (strong, nonatomic) AddressBookTopView *topView;
@property (strong, nonatomic) UITableView *tabel;

@end

@implementation AddressBookViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor colorWithRed:240 / 255.0 green:240 / 255.0 blue:240 / 255.0 alpha:1];
    [self creatTopView];
    [self.view addSubview:self.tabel];
}

- (void)creatTopView {
    _topView = [[AddressBookTopView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT / 4)];
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
    }
    return _tabel;
}

- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView {
    return 2;
}

- (NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section {
    return section == 0 ? 1 : 5;
}

- (CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath {
    return indexPath.section == 0 ? SCREEN_WIDTH / 3 + 20 : 120;
}

- (CGFloat)tableView:(UITableView *)tableView heightForHeaderInSection:(NSInteger)section {
    return 45;
}

- (CGFloat)tableView:(UITableView *)tableView heightForFooterInSection:(NSInteger)section {
    return 0.01;
}

- (UIView *)tableView:(UITableView *)tableView viewForHeaderInSection:(NSInteger)section {
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
    if (indexPath.section == 0) {
        static NSString *cellid = @"cellid";
        OneFriendTableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"cellid"];
        if (!cell) {
            cell = [[OneFriendTableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
        }
        return cell;
    }
    else {
        static NSString *cellid = @"cellid";
        FriendTableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"cellid"];
        if (!cell) {
            cell = [[FriendTableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
        }
        cell.headImg.image = [UIImage imageNamed:@"navLeftImg"];
        cell.nameLab.text = @"Crayon Shinchan";
        cell.messageLab.text = @"Nice to meet you!";
        return cell;
    }
}


@end
