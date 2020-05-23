//
//  PersonalInfoViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/16.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "PersonalInfoViewController.h"

@interface PersonalInfoViewController () <UITableViewDelegate, UITableViewDataSource, SettingImageViewDelegate, UINavigationControllerDelegate, UIImagePickerControllerDelegate, SetBottomSlideDelegate ,PersonalInfoTopViewDelegate>

{
    SetSlideCellTableViewCell *cell4;
}

@property (strong, nonatomic) PersonalInfoTopView *topView; // 头部视图
@property (strong, nonatomic) SettingImageView *setImgView; // 图片视图
@property (strong, nonatomic) UITableView *tabel;           // 表格

@property (strong, nonatomic) NSMutableArray *sectionArr;   // 标题
@property (strong, nonatomic) NSMutableArray *cellArr;      // 内容

@property (assign, nonatomic) NSInteger tag;                // 图片序号
@property (strong, nonatomic) NSMutableArray *imgArr;       // 存放图片信息

@property (strong, nonatomic) NSMutableArray *textArr;
@property (strong, nonatomic) NSString *minAge;
@property (strong, nonatomic) NSString *maxAge;

@property (strong, nonatomic) NSDictionary *dic;

@property (strong, nonatomic) NSMutableArray *imgUrlArr;

@end

@implementation PersonalInfoViewController



- (void)viewDidLoad {
    [super viewDidLoad];
    
    // 注册键盘方法
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardShowAndHide:) name:UIKeyboardWillShowNotification object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardShowAndHide:) name:UIKeyboardWillHideNotification object:nil];
    
    _minAge = @"16";
    _maxAge = @"30";
    _sectionArr = [[NSMutableArray alloc] initWithObjects:@"Name",@"Profession",@"Age",@"Bio",@"Seeking Age Range", nil];
    _textArr = [[NSMutableArray alloc] initWithObjects:@"",@"",@"",@"",@"",@"", nil];
    _cellArr = [[NSMutableArray alloc] initWithObjects:@"Name",@"Profession",@"Age",@"Bio",@"Min", @"max", nil];
    _imgArr = [[NSMutableArray alloc] initWithObjects:@"", @"", @"", nil];
    _imgUrlArr = [[NSMutableArray alloc] init];
    
    self.view.backgroundColor = [UIColor whiteColor];
    
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(changeText:) name:UITextFieldTextDidChangeNotification object:nil];
    
    [self requestDatas];
    [self creatTopView];
    [self creatSetImgView];
}

- (void)requestDatas {

    [HttpData requestInfoToken:[UserDefaults getUserDefaults:@"token"] success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"]integerValue] == 200) {
                self.dic = json[@"data"];
                [self.imgArr removeAllObjects];
                if (json[@"data"][@"photo"]) [self.imgArr addObject:json[@"data"][@"photo"]];
                if (json[@"data"][@"photo1"]) [self.imgArr addObject:json[@"data"][@"photo1"]];
                if (json[@"data"][@"photo2"]) [self.imgArr addObject:json[@"data"][@"photo2"]];
                if (json[@"data"][@"photo3"]) [self.imgArr addObject:json[@"data"][@"photo3"]];
                while (self.imgArr.count < 3) [self.imgArr addObject:@""];
                [self loadImages:self.imgArr];
                
                
                if (!(self.dic[@"name"] == nil || [self.dic[@"name"] isEqual:[NSNull null]])) self.textArr[0] = self.dic[@"name"];
                if (!(self.dic[@"work"] == nil || [self.dic[@"work"] isEqual:[NSNull null]])) self.textArr[1] = self.dic[@"work"];
                if (!(self.dic[@"age"] == nil || [self.dic[@"age"] isEqual:[NSNull null]])) self.textArr[2] = self.dic[@"age"];
                if (!(self.dic[@"bio"] == nil || [self.dic[@"bio"] isEqual:[NSNull null]])) self.textArr[3] = self.dic[@"bio"];
                
                if (!(self.dic[@"age_min"] == nil || [self.dic[@"age_min"] isEqual:[NSNull null]])) self.textArr[4] = self.dic[@"age_min"];
                else self.textArr[4] = @"16";
                if (!(self.dic[@"age_max"] == nil || [self.dic[@"age_max"] isEqual:[NSNull null]])) self.textArr[5] = self.dic[@"age_max"];
                else self.textArr[5] = @"100";
                self.minAge = self.textArr[4];
                self.maxAge = self.textArr[5];
                
                [self.view addSubview:self.tabel];
            }
            else {
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showErrorWithStatus:json[@"message"]];
            }
        }
    } failure:^(NSError * _Nonnull err) {
        [SVProgressHUD dismissWithDelay:1.5];
        [SVProgressHUD showErrorWithStatus:@"Error" ];
    }];
}

- (void)loadImages:(NSMutableArray*)arr {
    if (![arr[0] isEqualToString:@""]) [self.setImgView.image1 sd_setImageWithURL:[NSURL URLWithString:arr[0]] forState:UIControlStateNormal];
    if (![arr[1] isEqualToString:@""]) [self.setImgView.image2 sd_setImageWithURL:[NSURL URLWithString:arr[1]] forState:UIControlStateNormal];
    if (![arr[2] isEqualToString:@""]) [self.setImgView.image2 sd_setImageWithURL:[NSURL URLWithString:arr[2]] forState:UIControlStateNormal];
}

#pragma mark -- 顶部视图
- (void)creatTopView {
    _topView = [[PersonalInfoTopView alloc] initWithFrame: CGRectMake(0, 0, SCREEN_WIDTH, 60)];
    _topView.delegate = self;
    [self.view addSubview: _topView];
}

#pragma mark -- 图片视图
- (void)creatSetImgView {
    _setImgView = [[SettingImageView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_WIDTH)];
    _setImgView.delegate = self;
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
    if (indexPath.section != 4) {
        static NSString *cellid = @"cellid";
        UITableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:@"cellid"];
        if (!cell) {
            cell = [[UITableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
        }
        UITextField *textField = [[UITextField alloc] init];
        textField.clearButtonMode = UITextFieldViewModeWhileEditing;
        textField.font = Font(16);
        textField.tag = indexPath.section;
        textField.placeholder = _cellArr[indexPath.section];
        textField.text = _textArr[indexPath.section];
        [cell.contentView addSubview:textField];
        textField.sd_layout.leftSpaceToView(cell.contentView, 15).rightSpaceToView(cell.contentView, 15).topSpaceToView(cell.contentView, 0).bottomSpaceToView(cell.contentView, 0);
        
        return cell;
    }
    else {
        cell4 = [tableView dequeueReusableCellWithIdentifier:@"seller4"];
        if(!cell4) {
            cell4 = [[SetSlideCellTableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:@"seller4"];
        }
        cell4.delegate = self;
        [cell4 requestMinVale:_minAge maxValue:_maxAge];
        cell4.minLabel.text = [NSString stringWithFormat:@"Min: %@",_minAge];
        cell4.maxLabel.text = [NSString stringWithFormat:@"Max: %@",_maxAge];
        return cell4;
    }
}

#pragma mark -- 输入框内容改变
- (void)changeText:(NSNotification*)obj {
    UITextField *tf = obj.object;
    if (tf.tag == 0) [_textArr replaceObjectAtIndex:tf.tag withObject:tf.text];
    else if (tf.tag == 1) [_textArr replaceObjectAtIndex:tf.tag withObject:tf.text];
    else if (tf.tag == 2) [_textArr replaceObjectAtIndex:tf.tag withObject:tf.text];
    else if (tf.tag == 3) [_textArr replaceObjectAtIndex:tf.tag withObject:tf.text];
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
//    NSLog(@"%ld, %@",_tag, img);
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


# pragma mark - 键盘
- (void)keyboardShowAndHide:(NSNotification*)notification {
    NSDictionary* info = [notification userInfo];
    CGRect beginRect = [[info objectForKey:UIKeyboardFrameBeginUserInfoKey] CGRectValue];
    CGRect endRect = [[info objectForKey:UIKeyboardFrameEndUserInfoKey] CGRectValue];
    CGFloat change = (endRect.origin.y - beginRect.origin.y) / 2;
    [UIView animateWithDuration:0.25 animations:^{
        [self.view setFrame:CGRectMake(self.view.frame.origin.x, self.view.frame.origin.y + change, self.view.frame.size.width, self.view.frame.size.height)];
    }];
}

- (void)clickSlideValueTag:(NSInteger)tag withValue:(NSString *)value {
    switch (tag) {
        case 100:
            _minAge = value;
            self.textArr[4] = _minAge;
            cell4.minLabel.text = [NSString stringWithFormat:@"Min: %@",_minAge];
            break;
        case 101:
            _maxAge = value;
            self.textArr[5] = _maxAge;
            cell4.maxLabel.text = [NSString stringWithFormat:@"Max: %@",_maxAge];
            break;
            
        default:
            break;
    }
}

- (void)clickPersonalInfoTopViewButton:(NSInteger)tag {
    if (tag == 1) {
        [self dismissViewControllerAnimated:YES completion:nil];
    }
    else if (tag == 2) {
        UINavigationController *nvc = [[UINavigationController alloc]initWithRootViewController:[NSClassFromString(@"LoginViewController") new]];
        nvc.modalPresentationStyle = UIModalPresentationOverFullScreen;
        [nvc.navigationBar setHidden:YES];
        [self presentViewController:nvc animated:YES completion:nil];
    }
    else if (tag == 3) { // 保存
        if (_minAge > _maxAge) {
            [SVProgressHUD dismissWithDelay:1.5];
            [SVProgressHUD showErrorWithStatus:@"最小年龄不能大于最大年龄!"];
            return;
        }
        for (int i = 0; i < _textArr.count; i++) {
            if ([_textArr[i] isEqualToString:@""]) {
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showErrorWithStatus:@"请继续完善信息！"];
                return;
            }
        }
        NSLog(@"%@", _textArr);
        NSLog(@"%@", _imgArr);
        //上传图片
        for (int i = 0; i <_imgArr.count; i++) {
            
            if ([_imgArr[i] isKindOfClass:[UIImage class]]) {
                [HttpData requestImgUploadFile:_imgArr[i] success:^(id  _Nonnull json) {
                    if ([json isKindOfClass:[NSDictionary class]]) {
                        if ([json[@"status"]integerValue] == 200) {
                            [self.imgUrlArr addObject:json[@"data"]];
                        }
                    }
                    else {
                        [SVProgressHUD dismissWithDelay:1.5];
                        [SVProgressHUD showErrorWithStatus:json[@"message"]];
                        return;
                    }
                } failure:^(NSError * _Nonnull err) {
                    [SVProgressHUD dismissWithDelay:1.5];
                    [SVProgressHUD showErrorWithStatus:@"图片上传失败!"];
                    return;
                }];
            }
        }
//        NSLog(@"%@", self.imgUrlArr);
        // 保存信息
        [self saveInf];
    }
}

- (void)saveInf {
    [HttpData requestSaveOneSelfDataToken:[UserDefaults getUserDefaults:@"token"] photo:self.imgUrlArr.count > 0 ? self.imgUrlArr[0] : @"" photo1:self.imgUrlArr.count > 0 ? self.imgUrlArr[0] : @"" photo2:self.imgUrlArr.count > 1 ? self.imgUrlArr[1] : @"" photo3:self.imgUrlArr.count > 2 ? self.imgUrlArr[2] : @"" user_name:self.textArr[0] age:self.textArr[1] work:self.textArr[2] bio:self.textArr[3] age_min:self.textArr[4] age_max:self.textArr[5] success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"]integerValue] == 200) {
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showSuccessWithStatus:@"信息保存成功!"];
            }
        }
        else {
            [SVProgressHUD dismissWithDelay:1.5];
            [SVProgressHUD showErrorWithStatus:json[@"message"]];
        }
    } failure:^(NSError * _Nonnull err) {
        [SVProgressHUD dismissWithDelay:1.5];
        [SVProgressHUD showErrorWithStatus:@"信息保存失败!"];
    }];
}


@end
