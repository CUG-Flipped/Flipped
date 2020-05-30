//
//  ChatViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/24.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "ChatViewController.h"

@interface ChatViewController () <ChatTopViewDelegate, UITableViewDelegate, UITableViewDataSource>

@property (strong, nonatomic) ChatTopView *topView;

@property (strong, nonatomic) STInputBar *inputBar;

@property (strong, nonatomic) UITableView *table;

@property (strong, nonatomic) NSMutableArray *messageArr;

@end

@implementation ChatViewController

- (void)viewDidLoad {
    _messageArr = [[NSMutableArray alloc] init];
    [self getCHatDatas];
    [super viewDidLoad];
    [self creatTopView];
    [self.view addSubview:self.table];
    [self.view setBackgroundColor:[UIColor whiteColor]];
    [self chreatKeyBoard];
    
}

- (void)creatTopView {
    _topView = [[ChatTopView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT / 7)];
    _topView.backgroundColor = [UIColor colorWithRed:240/255.0 green:240/255.0 blue:240/255.0 alpha:0.5];
    _topView.delegate = self;
    [_topView.headImg sd_setImageWithURL:url placeholderImage:[UIImage imageNamed:@"navLeftImg"]];
    _topView.nameLabel.text = name;
    [self.view addSubview:self.topView];
}

- (void)chatTopViewButtonClick:(NSInteger)tag {
    if (tag == 1) {
        [self.navigationController popViewControllerAnimated:YES];
    }
    if (tag == 2) {
        
    }
}

- (void)chreatKeyBoard {
    _inputBar = [STInputBar inputBar];
    _inputBar.center = CGPointMake(self.view.bounds.size.width / 2, self.view.bounds.size.height - _inputBar.bounds.size.height / 2);
    [_inputBar setFitWhenKeyboardShowOrHide:YES];
    _inputBar.placeHolder = @"Please enter here...";
    [self.view addSubview:_inputBar];
    
    kWeakSelf(self);
    _inputBar.sendDidClickedHandler = ^(NSString *string) {
        if ([string isEqualToString:@""] || string.length == 0) {
            return;
        }
        else {
            weakself.inputBar.textView.text = @"";
            
            NSDictionary *dic = @{
                @"content" : string,
                @"type": @"0"
            };
            
            [weakself.messageArr addObject:dic];
            [weakself.table reloadData];
            
            [weakself quitEdit];
        }
    };
}

- (void)quitEdit {
    [_inputBar resignFirstResponder];
}

- (UITableView*)table {
    if (!_table) {
        _table = [[UITableView alloc] initWithFrame:CGRectMake(0, self.topView.bounds.size.height, SCREEN_WIDTH, SCREEN_HEIGHT - self.topView.bounds.size.height - _inputBar.bounds.size.height) style:UITableViewStylePlain];
        _table.delegate = self;
        _table.dataSource = self;
//        _table.backgroundColor = [UIColor colorWithRed:245/255.0 green:245/255.0 blue:245/255.0 alpha:1];
        _table.separatorStyle = UITableViewCellSeparatorStyleNone;
    }
    return _table;
}
- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView {
    return 1;
}

- (NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section {
    return _messageArr.count;
}
- (CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath {
    
    NSDictionary *dic = @{ NSFontAttributeName : [UIFont systemFontOfSize:16] };
    CGSize maxSize = CGSizeMake(SCREEN_WIDTH /  3 * 2, MAXFLOAT);
    NSStringDrawingOptions option = NSStringDrawingUsesLineFragmentOrigin | NSStringDrawingUsesFontLeading;
    CGSize size = [_messageArr[indexPath.row][@"content"] boundingRectWithSize:maxSize options:option attributes:dic context:nil].size;
    
    return ceilf(size.height + 30); // 向上取整
//    NSStringDrawingUsesLineFragmentOrigin 整个文本将以每行组成的矩形为单位计算整个文本的尺寸
//    NSStringDrawingUsesFontLeading 使用字体的行间距来计算文本占用的范围，即每一行的底部到下一行的底部的距离计算
}
- (UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath {
    static NSString *cellid = @"cellid";
    ChatTableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:cellid];
    if (!cell) {
        cell = [[ChatTableViewCell alloc] initWithStyle:UITableViewCellStyleDefault reuseIdentifier:cellid];
    }
    cell.messageLabel.text = _messageArr[indexPath.row][@"content"];
    
    NSDictionary *dic = @{ NSFontAttributeName : [UIFont systemFontOfSize:16] };
    CGSize maxSize = CGSizeMake(SCREEN_WIDTH / 3 * 2, MAXFLOAT);
    NSStringDrawingOptions option = NSStringDrawingUsesLineFragmentOrigin | NSStringDrawingUsesFontLeading;
    CGSize size = [_messageArr[indexPath.row][@"content"] boundingRectWithSize:maxSize options:option attributes:dic context:nil].size;
//    NSInteger hightC = size.height;
    NSInteger widthC = size.width;
//    CGSize sizeTpFit = [_messageArr[indexPath.row][@"content"] sizeWithFont:Font(16) constrainedToSize:CGSizeMake(SCREEN_WIDTH / 3 * 2, hightC - 10) lineBreakMode:NSLineBreakByWordWrapping];
    
    if ([_messageArr[indexPath.row][@"type"]integerValue] == 0) {
        cell.backButton.backgroundColor = [UIColor colorWithRed:51 / 255.0 green:153 / 255.0 blue:255 / 255.0 alpha:1];
        cell.messageLabel.textColor = [UIColor whiteColor];
        cell.backButton.sd_layout.rightSpaceToView(cell.contentView, 15).topSpaceToView(cell.contentView, 5).widthIs(widthC + 30 ).bottomSpaceToView(cell.contentView, 5);
    }
    else {
        cell.backButton.backgroundColor = [UIColor colorWithRed:210 / 255.0 green:210 / 255.0 blue:210 / 255.0 alpha:1];
        cell.messageLabel.textColor = [UIColor blackColor];
        cell.backButton.sd_layout.leftSpaceToView(cell.contentView, 15).topSpaceToView(cell.contentView, 5).widthIs(widthC + 30 ).bottomSpaceToView(cell.contentView, 5);
    }
    
    cell.selectionStyle = UITableViewCellSelectionStyleNone;
    return cell;
}

- (void)getCHatDatas {
    NSDictionary *dic1 = @{
        @"content" : @"Hello! We are friends now, let's chat!",
        @"type": @"0"
    };
    NSDictionary *dic2 = @{
        @"content" : @"Hi! Nice to meet you!",
        @"type": @"1"
    };
    NSDictionary *dic3 = @{
        @"content" : @"Where do you go to college?",
        @"type": @"1"
    };
    NSDictionary *dic4 = @{
        @"content" : @"I am in the School of Geography and Information Engineering of China University of Geosciences (Wuhan)! Software engineering major!",
        @"type": @"0"
    };
    NSDictionary *dic5 = @{
        @"content" : @"Wow! Great!!",
        @"type": @"1"
    };
    [_messageArr addObject:dic1];
    [_messageArr addObject:dic2];
    [_messageArr addObject:dic3];
    [_messageArr addObject:dic4];
    [_messageArr addObject:dic5];
}

@end
