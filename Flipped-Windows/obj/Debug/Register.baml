   M S B A M L   `   `   `  ����H  DFlipped-Win10, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null-clr-namespace:Flipped_Win10Flipped_Win10  5   	       Flipped_Win10.Register   5   	   W SPresentationCore, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35\ XPresentationFramework, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35R NWindowsBase, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35D 9http://schemas.microsoft.com/winfx/2006/xaml/presentation    5   	   R NSystem.Xaml, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089:x,http://schemas.microsoft.com/winfx/2006/xaml     5   	   @d2http://schemas.microsoft.com/expression/blend/2008     5   	   Jmc;http://schemas.openxmlformats.org/markup-compatibility/2006     5   	   -localclr-namespace:Flipped_Win10     5   	   V RMaterialDesignThemes.Wpf, Version=3.1.2.1465, Culture=neutral, PublicKeyToken=nullFmaterialDesign1http://materialdesigninxaml.net/winfx/xaml/themes  5   	   -   5          �Name$  window��5   	    � Title$ Register��5	   	   $	��600��6   $	��370��6'    � WindowStyle$
 None=�63    � AllowsTransparency .  � WindowStartupLocation$ CenterScreen=�6`    � 
ResizeMode$ NoResize=�6�    � 	Resources 5
         registerWin_LoadStory �� TargetProperty �� 
TargetName9  3System.Windows.Media.Animation.EasingDoubleKeyFrame	 O� KeyTime
 O� Value��%�   &	        �� 5   
   .��5      U� $, &(UIElement.Effect).(BlurEffect.Radius)�6,   $	 box��6o   .\�5        $	 0��6'   $
 20Q�63     5      $	 0:0:0.5��6'   $
 0Q�69   5      5      5       � Triggers 5      :� 5   
    :� RoutedEvent$ FrameworkElement.Loaded��6   .[�5      �� #��[  6   5      '  !System.Windows.Shell.WindowChrome   WindowChrome 5        5   
      GlassFrameThickness$ 0q�6      ResizeBorderThickness$ 10q�6/      CaptionHeight$ 50Q�6J   ��5      �� 5   
    �� 
TargetType" �6   .R�5      �  �Name$ backGrid��6   .+�5      ��  ��Name$ contentBorder��6   $��10q�60    �� CornerRadius$ 10��6<   ���#fefefc J� Effect 5      3  -System.Windows.Media.Effects.DropShadowEffect  5         
BlurRadius$ 6Q�6+      ShadowDepth$ 0Q�6:      Color$ #B7B7B7��6J    J� OpacityMask 5!      !� 5"       !� Visual 5#      �� 5$   "   
���   � J� SnapsToDevicePixels .  5&   )   �� *CornerRadius+ �� RelativeSource ��  �� AncestorType$ BorderL���5'   )   �� *ActualWidth+ �� $ BorderL���5(   )   �� *ActualHeight+ �� $ BorderL�5*      ._�5,      �� .u�5-      �� 5.      5/      50       �� Triggers 51      P� 52       P� Property �@�IsActive  P� Value  . . �53      �� ! �� 
TargetName$! contentBorder��6   " �� Property" �@J�Effect# �� Value# 54        55      $ 9Q�6/   $ 0Q�6>   $ #B2B2B2��6N   57      58      5:      .��5<      � }�. .+�5=   
   � -   5        $	 box��5=      ��5>      �� 5?      $ �� ImageSource$$ /Images/register_bg.png��6    5A      -  'System.Windows.Media.Effects.BlurEffect  5B      %   
KernelType$% Gaussian=�6   &   Radius$& 0Q�63   .+�5D      �� -   5        $
��Left=�5D      $	��Top=�6K     WindowCloseButtonStyle#��[ 6c   $��30��6�   $��60��6�   '   IsHitTestVisibleInChrome' . $��	292,0,0,0q�6�   �� 5E      $
��Left=�6   $��45��64   $��
55,150,0,0q�6@   $	��Top=�6T   $	��245��6l   ( �� Orientation$( 
Horizontal=�6x   .+�5F      '  !MaterialDesignThemes.Wpf.PackIcon  )   Kind$) AccountCircle=�6*   $��20��6?   $��20��6K   $
��Left=�6V   $��Bottom=�6q   $��0,0,5,5q�6�   �� 5G      -   5        * ��Name$* userNameBox��5G      )  #MaterialDesignThemes.Wpf.HintAssist+   Hint$+ UserName��6/    % !MaterialDesignFloatingHintTextBox#��[ 6Y   .  (MaterialDesignThemes.Wpf.TextFieldAssist,   HasClearButton, . $��15(�5H      - �� TextWrapping$
- Wrap=�6)   $	��200��6=   $��5,0,0,0q�6I   5I      �� 5J      $
��Left=�6   $��45��64   $��
55,200,0,0q�6@   $	��Top=�6T   $	��245��6l   $( 
Horizontal=�6x   .+�5K        $) Onepassword=�6*   $��20��6=   $��20��6I   $
��Left=�6T   $��Bottom=�6o   $��0,0,5,5q�6�   G� 5L      -   5        . G�Name$. pwdBox��5L      $+ Password��6.    ) %MaterialDesignFloatingHintPasswordBox#��[ 6X   /   
Foreground
/ ��  �0   UnderlineBrush
0 ��  �, . $��15(�5N      $	��200��6-   1 G� 	MaxLength$1 20��69   $��5,0,0,0q�6H   5O      �� 5P      $
��Left=�6   $��45��64   $��
55,250,0,0q�6@   $	��Top=�6T   $	��245��6l   $( 
Horizontal=�6x   .+�5Q        $) EmailOutline=�6*   $��20��6>   $��20��6J   $
��Left=�6U   $��Bottom=�6p   $��0,0,5,5q�6�   �� 5R      -   5        $* emailBox��5R      $+ Email��6,   #��[ 6S   
/ ��  �
0 ��  �, . $��15(�5T      $	��200��6-   $��5,0,0,0q�69   5U      �� 5V      $
��Left=�6   $��45��64   $��
55,300,0,0q�6@   $	��Top=�6T   $	��245��6l   $( 
Horizontal=�6x   .+�5W        $) CardAccountDetails=�6*   $��20��6D   $��20��6P   $
��Left=�6[   $��Bottom=�6v   $��0,0,5,5q�6�   �� 5X      -   5        2 ��Name$2 userTypeBox��5X      $+ userType��60   $��Center=�6Z   $��90��6|   $��45��6�   $��15(�6�   $��5,0,0,0q�6�   .@�5Y      �� $��普通用户��6#   �� 5Z      $��	vip用户��6#   �� 5[      $��	管理员��6#   5\      5]      �� 5^      $
��Left=�6   $	��125��64   $��
55,350,0,0q�6A   $	��Top=�6U   $	��245��6m   $( 
Horizontal=�6y   .+�5_        $) Upload=�6*   $��20��68   $��20��6C   $��Center=�6O   �� 5a      -   5          MaterialDesignFlatButton#��[ 5a      
���   �$��Upload��6{   $��80��6�   �� 5b      -	   5        3 ��Name$3 avatar��5b      $	��125��6(   $��Right=�64   $��10,0,0,0q�6P   $��Uniform=�6b   5c      �� 5d      -
   5        4 ��Name$4 registerBtn��5d      $��Register��6D   $
��Left=�6W   $��125,515,0,0q�6r   $	��Top=�6�   $	��100��6�   5 J� RenderTransformOrigin$5 0,0.469&�6�   +  %MaterialDesignThemes.Wpf.ShadowAssist6   ShadowDepth$6 Depth4=�6�   �� 5e      $��Register��6   $��Center=�6'   $��Center=�6K   $��50(�6m   
�������$
��Bold"�6�   $��Gabriola)�6�   $
��Left=�6�   $��	50,45,0,0q�6�   $	��Top=�6�   $	��100��6  $	��250��6  '	  !MaterialDesignThemes.Wpf.Snackbar	  5g      -   5        7 	 Name$7 snackBar��5g   &   8 	  IsActive8 .  9 	  Message.9 5h      .
  (MaterialDesignThemes.Wpf.SnackbarMessage
  -   5        : 
 Name$: snackBarMsg��5h   1   $)��#Error! Please Check your infomation��6F   ; 
  ActionContent$; OK��6t   6�   5i      5j      5k      5m      