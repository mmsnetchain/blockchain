package config

const (
	MSGID_search_node           = 135
	MSGID_search_node_recv      = 136
	MSGID_checkNodeOnline       = 110
	MSGID_checkNodeOnline_recv  = 111
	MSGID_TextMsg               = 112
	MSGID_getNearSuperIP        = 113
	MSGID_getNearSuperIP_recv   = 114
	MSGID_multicast_online_recv = 122
	MSGID_ask_close_conn        = 127
	MSGID_ask_close_conn_recv   = 128
	MSGID_TextMsg_recv          = 129

	MSGID_SearchAddr                = 130
	MSGID_SearchAddr_recv           = 131
	MSGID_security_create_pipe      = 132
	MSGID_security_create_pipe_recv = 133
	MSGID_security_pipe_error       = 134

	MSGID_register_name        = 102
	MSGID_register_name_recv   = 103
	MSGID_build_name           = 104
	MSGID_build_name_recv      = 105
	MSGID_check_temp_name      = 106
	MSGID_check_temp_name_recv = 107

	MSGID_find_name                = 115
	MSGID_find_name_recv           = 116
	MSGID_name_add_address_recv    = 117
	MSGID_name_sync_multicast_recv = 118

	MSGID_ROOT_register_name    = 123
	MSGID_ROOT_RECV_create_name = 124

	MSGID_key_sync_multicast_recv = 119
	MSGID_key_find_keyname        = 120
	MSGID_key_find_keyname_recv   = 121
	MSGID_ROOT_RECV_save_key_name = 125

	MSGID_http_request  = 126
	MSGID_http_response = 127

	MSGID_multicast_vote_recv      = 200
	MSGID_multicast_blockhead      = 201
	MSGID_heightBlock              = 202
	MSGID_heightBlock_recv         = 203
	MSGID_getStartBlockHead        = 204
	MSGID_getStartBlockHead_recv   = 205
	MSGID_getBlockHeadVO           = 206
	MSGID_getBlockHeadVO_recv      = 207
	MSGID_multicast_transaction    = 208
	MSGID_getUnconfirmedBlock      = 209
	MSGID_getUnconfirmedBlock_recv = 210

	MSGID_multicast_return        = 211
	MSGID_getblockforwitness      = 212
	MSGID_getblockforwitness_recv = 213

	MSGID_getDBKey_one      = 214
	MSGID_getDBKey_one_recv = 215

	MSGID_multicast_find_witness      = 216
	MSGID_multicast_find_witness_recv = 217

	MSGID_getBlockLastCurrent      = 218
	MSGID_getBlockLastCurrent_recv = 219

	MSGID_multicast_witness_blockhead          = 220
	MSGID_multicast_witness_blockhead_recv     = 221
	MSGID_multicast_witness_blockhead_get      = 222
	MSGID_multicast_witness_blockhead_get_recv = 223

	MSGID_uniformity_multicast_witness_blockhead         = 224
	MSGID_uniformity_multicast_witness_blockhead_recv    = 225
	MSGID_uniformity_multicast_witness_block_get         = 226
	MSGID_uniformity_multicast_witness_block_get_recv    = 227
	MSGID_uniformity_multicast_witness_block_import      = 228
	MSGID_uniformity_multicast_witness_block_import_recv = 229

	MSGID_sharebox_addFileShare                     = 300
	MSGID_sharebox_addFileShare_recv                = 301
	MSGID_sharebox_findFileinfo                     = 302
	MSGID_sharebox_findFileinfo_recv                = 303
	MSGID_sharebox_getFilesize                      = 304
	MSGID_sharebox_getFilesize_recv                 = 305
	MSGID_sharebox_downloadFileChunk                = 306
	MSGID_sharebox_downloadFileChunk_recv           = 307
	MSGID_sharebox_getUploadinfo                    = 308
	MSGID_sharebox_getUploadinfo_recv               = 309
	MSGID_sharebox_getNodeWalletReceiptAddress      = 315
	MSGID_sharebox_getNodeWalletReceiptAddress_recv = 316
	MSGID_sharebox_getsharefolderlist               = 317
	MSGID_sharebox_getsharefolderlist_recv          = 318

	MSGID_store_addFileShare      = 400
	MSGID_store_addFileShare_recv = 401
	MSGID_store_findFileinfo      = 402
	MSGID_store_findFileinfo_recv = 403

	MSGID_store_downloadFileChunk      = 406
	MSGID_store_downloadFileChunk_recv = 407

	MSGID_store_syncFileInfo                     = 410
	MSGID_store_syncFileInfo_recv                = 412
	MSGID_store_getfourNodeinfo                  = 413
	MSGID_store_getfourNodeinfo_recv             = 414
	MSGID_store_getNodeWalletReceiptAddress      = 415
	MSGID_store_getNodeWalletReceiptAddress_recv = 416

	MSGID_store_online_heartbeat      = 417
	MSGID_store_online_heartbeat_recv = 418
	MSGID_store_getFileindexList      = 419
	MSGID_store_getFileindexList_recv = 420
	MSGID_store_addFileOwner          = 421
	MSGID_store_addFileOwner_recv     = 422

	MSGID_im_addfriend        = 500
	MSGID_im_addfriend_recv   = 501
	MSGID_im_agreefriend      = 502
	MSGID_im_agreefriend_recv = 503
	MSGID_im_file             = 504
	MSGID_im_file_recv        = 505
	MSGID_im_property         = 506
	MSGID_im_property_recv    = 507
	MSGID_im_addr             = 508
	MSGID_im_addr_recv        = 509
	MSGID_im_pay              = 5010
	MSGID_im_pay_recv         = 5011

	MSGID_vnode_getstate            = 600
	MSGID_vnode_getstate_recv       = 601
	MSGID_vnode_getNearSuperIP      = 602
	MSGID_vnode_getNearSuperIP_recv = 603

	MSGID_imnew_upload_userinfo   = 700
	MSGID_imnew_get_userinfo      = 701
	MSGID_imnew_get_userinfo_recv = 702

	MSGID_imnew_agreefriend_recv = 703
	MSGID_imnew_file             = 704
	MSGID_imnew_file_recv        = 705
	MSGID_imnew_property         = 706
	MSGID_imnew_property_recv    = 707
	MSGID_imnew_addr             = 708
	MSGID_imnew_addr_recv        = 709
)
