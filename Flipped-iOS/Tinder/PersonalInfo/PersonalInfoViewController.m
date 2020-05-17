//
//  PersonalInfoViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/16.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "PersonalInfoViewController.h"

@interface PersonalInfoViewController () <UITableViewDelegate, UITableViewDataSource, SettingImageViewDelegate, UINavigationControllerDelegate, UIImagePickerControllerDelegate>

@property (strong, nonatomic) PersonalInfoTopView *topView;
@property (strong, nonatomic) UITableView *tabel;

@property (strong, nonatomic) NSMutableArray *sectionArr;
@property (strong, nonatomic) NSMutableArray *cellArr;

@property (strong, nonatomic) SettingImageView *setImgView;

@property (assign, nonatomic) NSInteger tag;
@property (strong, nonatomic) NSMutableArray *imgArr;

@end

@implementation PersonalInfoViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
    
    _sectionArr = [[NSMutableArray alloc] initWithObjects:@"Name",@"Profession",@"Age",@"Bio",@"Seeking Age Range", nil];
    _cellArr = [[NSMutableArray alloc] initWithObjects:@"Crayon Shinchan",@"Student",@"5",@"Nice to meet you!", nil];
    _imgArr = [[NSMutableArray alloc] initWithObjects:@"background_1", @"background_2", @"background_3", nil];
    
    self.view.backgroundColor = [UIColor whiteColor];
    [self creatTopView];
    [self creatSetImgView];
    [self.view addSubview:self.tabel];
}

- (void)creatTopView {
    _topView = [[PersonalInfoTopView alloc] initWithFrame: CGRectMake(0, 0, SCREEN_WIDTH, 60)];
    
    [self.view addSubview: _topView];
}

#pragma mark -- 图片视图
- (void)creatSetImgView {
    _setImgView = [[SettingImageView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_WIDTH)];
    _setImgView.delegate = self;
//    [self.view addSubview:self.setImgView];
//    _setImgView.sd_layout.leftSpaceToView(self.view, 0).topSpaceToView(_topView, 0).rightSpaceToView(self, 0).heightIs(SCREEN_WIDTH).widthIs(SCREEN_WIDTH);
}

#pragma mark -- 表格
- (UITableView*)tabel {
    if (!_tabel) {
        _tabel = [[UITableView alloc] initWithFrame:CGRectMake(0, self.topView.frame.size.height, SCREEN_WIDTH, SCREEN_HEIGHT - self.topView.frame.size.height) style:UITableViewStyleGrouped];
        _tabel.delegate = self;
        _tabel.dataSource = self;
        _tabel.bounces = NO; // 弹簧效果
        _tabel.separatorStyle = UITableViewCellSeparatorStyleNone;
    }
    return _tabel;
}

// 分区
- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView {
    return _sectionArr.count;
}
// 每个分区几行
- (NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section {
    return 1;
}
//  每一行多高
- (CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath {
    return indexPath.section == _sectionArr.count - 1 ? 150 : 50;
}
// 灰色条高
- (CGFloat)tableView:(UITableView *)tableView heightForHeaderInSection:(NSInteger)section {
    return section == 0 ? SCREEN_WIDTH : 50;
}
// 灰色条底
-(CGFloat)tableView:(UITableView *)tableView heightForFooterInSection:(NSInteger)section {
    return 0.01;
}
// 头部视图
- (UIView *)tableView:(UITableView *)tableView viewForHeaderInSection:(NSInteger)section {
    if (section == 0) return _setImgView;
    UIView *view = [[UIView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, 50)];
    UILabel *label = [[UILabel alloc] initWithFrame:CGRectMake(15, 0, SCREEN_WIDTH - 30, 50)];
    label.text = _sectionArr[section];
    label.font = BoldFont(20);
    [view addSubview:label];
    return view;
}
// 底视图
- (UIView *)tableView:(UITableView *)tableView viewForFooterInSection:(NSInteger)section {
    return nil;
}

// 每一行内容
- (UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath {
    static NSString *cellid = @"cellid";
    UITableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"cellid"];
    if (!cell) {
        cell = [[UITableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
    }
    cell.textLabel.text = @"1";
    return cell;
}

- (void)scrollViewDidScroll:(UIScrollView *)scrollView {
    [self.view endEditing:YES];
    if (scrollView.contentOffset.y >= 60) {
        _topView.titleLabel.hidden = NO;
        _topView.backgroundColor = [UIColor colorWithRed:221 / 255.0 green:221 / 255.0 blue:221 / 255.0 alpha:0.5];
    }
    else {
        _topView.titleLabel.hidden = YES;
        _topView.backgroundColor = [UIColor whiteColor];
    }
}

- (void)clickImageButton:(NSInteger)tag {
//    NSLog(@"%ld",tag);
    _tag = tag;
    UIImagePickerController *pick = [[UIImagePickerController alloc] init];
    pick.allowsEditing = YES;
    pick.delegate = self;
    pick.sourceType = UIImagePickerControllerSourceTypePhotoLibrary;
    [self presentViewController:pick animated:YES completion:nil];
}

- (void)imagePickerController:(UIImagePickerController *)picker didFinishPickingMediaWithInfo:(NSDictionary<UIImagePickerControllerInfoKey,id> *)info {
    [picker dismissViewControllerAnimated:YES completion:nil];
    UIImage *img = [info objectForKey:UIImagePickerControllerEditedImage];
    if (!img) return;
    NSLog(@"%ld, %@",_tag, img);
    [_imgArr replaceObjectAtIndex:_tag - 1 withObject: img];
    switch (_tag) {
        case 1:
            [_setImgView.image1 setImage:img forState:UIControlStateNormal];
            break;
        case 2:
            [_setImgView.image2 setImage:img forState:UIControlStateNormal];
            break;
        case 3:
            [_setImgView.image3 setImage:img forState:UIControlStateNormal];
            break;
        default:
            break;
    }
}

@end
